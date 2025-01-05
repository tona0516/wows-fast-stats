package service

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"sync"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/repository"
	"wfs/backend/yamibuka"

	"github.com/abadojack/whatlanggo"
	"github.com/morikuni/failure"
)

type Battle struct {
	wargaming      repository.WargamingInterface
	uwargaming     repository.UnofficialWargamingInterface
	numbers        repository.NumbersInterface
	unregistered   repository.UnregisteredInterface
	localFile      repository.LocalFileInterface
	storage        repository.StorageInterface
	logger         repository.LoggerInterface
	eventsEmitFunc eventEmitFunc

	isFirstBattle                      bool
	isNotifyExpectedStatsUnavaillalble bool
	warship                            data.Warships
	allExpectedStats                   data.ExpectedStats
	battleArenas                       data.WGBattleArenas
	battleTypes                        data.WGBattleTypes
}

func NewBattle(
	wargaming repository.WargamingInterface,
	uwargaming repository.UnofficialWargamingInterface,
	localFile repository.LocalFileInterface,
	numbers repository.NumbersInterface,
	unregistered repository.UnregisteredInterface,
	storage repository.StorageInterface,
	logger repository.LoggerInterface,
	eventsEmitFunc eventEmitFunc,
) *Battle {
	return &Battle{
		wargaming:                          wargaming,
		uwargaming:                         uwargaming,
		localFile:                          localFile,
		numbers:                            numbers,
		unregistered:                       unregistered,
		storage:                            storage,
		logger:                             logger,
		eventsEmitFunc:                     eventsEmitFunc,
		isFirstBattle:                      true,
		isNotifyExpectedStatsUnavaillalble: false,
	}
}

func (b *Battle) Get(appCtx context.Context, userConfig data.UserConfigV2) (data.Battle, error) {
	var result data.Battle

	// Fetch on-memory stored data
	warshipResult := make(chan data.Result[data.Warships])
	allExpectedStatsResult := make(chan data.Result[data.ExpectedStats])
	battleArenasResult := make(chan data.Result[data.WGBattleArenas])
	battleTypesResult := make(chan data.Result[data.WGBattleTypes])
	if b.isFirstBattle {
		go b.fetchWarships(warshipResult)
		go b.fetchExpectedStats(allExpectedStatsResult)
		go b.fetchBattleArenas(battleArenasResult)
		go b.fetchBattleTypes(battleTypesResult)
	}

	// Get tempArenaInfo.json
	tempArenaInfo, err := b.getTempArenaInfo(userConfig)
	if err != nil {
		return result, err
	}

	// persist own ign for reporting
	_ = b.storage.WriteOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	// Get Account ID list
	accountList, err := b.wargaming.AccountList(tempArenaInfo.AccountNames())
	if err != nil {
		return result, err
	}
	accountIDs := accountList.AccountIDs()

	// Fetch each stats
	accountInfoResult := make(chan data.Result[data.WGAccountInfo])
	shipStatsResult := make(chan data.Result[data.AllPlayerShipsStats])
	clanResult := make(chan data.Result[data.Clans])
	go b.fetchAccountInfo(accountIDs, accountInfoResult)
	go b.fetchAllPlayerShipsStats(accountIDs, shipStatsResult)
	go b.fetchClan(accountIDs, clanResult)

	errs := make([]error, 0)

	if b.isFirstBattle {
		warship := <-warshipResult
		b.warship = warship.Value
		errs = append(errs, warship.Error)

		expectedStats := <-allExpectedStatsResult
		b.allExpectedStats = expectedStats.Value
		errs = append(errs, expectedStats.Error)

		battleArenas := <-battleArenasResult
		b.battleArenas = battleArenas.Value
		errs = append(errs, battleArenas.Error)

		battleTypes := <-battleTypesResult
		b.battleTypes = battleTypes.Value
		errs = append(errs, battleTypes.Error)
	}

	accountInfo := <-accountInfoResult
	errs = append(errs, accountInfo.Error)

	shipStats := <-shipStatsResult
	errs = append(errs, shipStats.Error)

	clan := <-clanResult
	errs = append(errs, clan.Error)

	for _, err := range errs {
		if err != nil {
			if failure.Is(err, apperr.ExpectedStatsUnavaillalble) && !b.isNotifyExpectedStatsUnavaillalble {
				b.eventsEmitFunc(appCtx, EventErr, apperr.ExpectedStatsUnavaillalble.ErrorCode())
				b.isNotifyExpectedStatsUnavaillalble = true
				continue
			}
			return result, err
		}
	}

	result = b.compose(
		tempArenaInfo,
		accountInfo.Value,
		accountList,
		clan.Value,
		shipStats.Value,
		b.warship,
		b.allExpectedStats,
		b.battleArenas,
		b.battleTypes,
	)

	b.isFirstBattle = false

	return result, nil
}

func (b *Battle) getTempArenaInfo(userConfig data.UserConfigV2) (data.TempArenaInfo, error) {
	tempArenaInfo, err := b.localFile.TempArenaInfo(userConfig.InstallPath)
	if err != nil {
		return tempArenaInfo, err
	}

	if userConfig.SaveTempArenaInfo {
		if err := b.localFile.SaveTempArenaInfo(tempArenaInfo); err != nil {
			return tempArenaInfo, err
		}
	}

	return tempArenaInfo, nil
}

func (b *Battle) fetchWarships(channel chan data.Result[data.Warships]) {
	warships := make(data.Warships)
	var result data.Result[data.Warships]

	var mu sync.Mutex

	fetch := func(page int) (int, error) {
		res, pageTotal, err := b.wargaming.EncycShips(page)
		if err != nil {
			return 0, err
		}

		for shipID, warship := range res {
			mu.Lock()
			warships[shipID] = data.Warship{
				Name:      warship.Name,
				Tier:      warship.Tier,
				Type:      data.NewShipType(warship.Type),
				Nation:    data.Nation(warship.Nation),
				IsPremium: warship.IsPremium,
			}
			mu.Unlock()
		}
		return pageTotal, nil
	}

	first := 1
	pageTotal, err := fetch(first)
	if err != nil {
		result.Error = err
		channel <- result
		return
	}

	pages := makeRange(first+1, pageTotal+1)
	err = doParallel(pages, func(page int) error {
		_, err := fetch(page)
		return err
	})
	if err != nil {
		result.Error = err
		channel <- result
		return
	}

	unregisteredShipInfo, err := b.unregistered.Warship()
	if err != nil {
		result.Error = err
		channel <- result
		return
	}
	for k, v := range unregisteredShipInfo {
		if _, ok := warships[k]; !ok {
			warships[k] = v
		}
	}

	result.Value = warships
	channel <- result
}

func (b *Battle) fetchExpectedStats(channel chan data.Result[data.ExpectedStats]) {
	var result data.Result[data.ExpectedStats]

	// 最新の予測成績を取得
	expectedStats, errFetch := b.numbers.ExpectedStats()
	if errFetch == nil {
		_ = b.storage.WriteExpectedStats(expectedStats)
		result.Value = expectedStats
		channel <- result
		return
	}

	// 取得できない場合、キャッシュを利用する
	expectedStats, errCache := b.storage.ExpectedStats()
	if errCache == nil {
		result.Value = expectedStats
		channel <- result
		return
	}

	result.Error = failure.New(apperr.ExpectedStatsUnavaillalble, failure.Context{
		"err_fetch": errFetch.Error(),
		"err_cache": errCache.Error(),
	})
	channel <- result
}

func (b *Battle) fetchBattleArenas(channel chan data.Result[data.WGBattleArenas]) {
	battleArenas, err := b.wargaming.BattleArenas()
	channel <- data.Result[data.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (b *Battle) fetchBattleTypes(channel chan data.Result[data.WGBattleTypes]) {
	battleTypes, err := b.wargaming.BattleTypes()
	channel <- data.Result[data.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (b *Battle) fetchAccountInfo(accountIDs []int, channel chan data.Result[data.WGAccountInfo]) {
	accountInfo, err := b.wargaming.AccountInfo(accountIDs)
	channel <- data.Result[data.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (b *Battle) fetchAllPlayerShipsStats(
	accountIDs []int,
	channel chan data.Result[data.AllPlayerShipsStats],
) {
	shipStatsMap := make(data.AllPlayerShipsStats)
	var mu sync.Mutex
	err := doParallel(accountIDs, func(accountID int) error {
		shipStats, err := b.wargaming.ShipsStats(accountID)
		if err != nil {
			return err
		}

		mu.Lock()
		shipStatsMap[accountID] = shipStats
		mu.Unlock()

		return nil
	})

	channel <- data.Result[data.AllPlayerShipsStats]{Value: shipStatsMap, Error: err}
}

func (b *Battle) fetchClan(accountIDs []int, channel chan data.Result[data.Clans]) {
	var result data.Result[data.Clans]

	clansAccountInfo, err := b.wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		result.Error = err
		channel <- result
		return
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := b.wargaming.ClansInfo(clanIDs)
	if err != nil {
		result.Error = err
		channel <- result
		return
	}

	clanInfoArray := clansInfo.ToArray()
	colorMap := b.fetchClanColor(clanInfoArray)
	languageMap := b.fetchClanLanguage(clanInfoArray)

	clans := make(data.Clans)
	for _, accountID := range accountIDs {
		clanID := clansAccountInfo[accountID].ClanID
		clanTag := clansInfo[clanID].Tag
		hexColor := colorMap[clanTag]
		language := languageMap[clanTag]

		clans[accountID] = data.Clan{Tag: clanTag, ID: clanID, HexColor: hexColor, Language: language}
	}

	result.Value = clans
	channel <- result
}

func (b *Battle) fetchClanColor(clanInfoArray []data.WGClansInfoData) map[string]string {
	result := make(map[string]string)

	var mu sync.Mutex
	err := doParallel(clanInfoArray, func(clan data.WGClansInfoData) error {
		autocomplete, err := b.uwargaming.ClansAutoComplete(clan.Tag)
		if err != nil {
			return err
		}

		hexColor := autocomplete.HexColor(clan.Tag)
		if hexColor != "" {
			mu.Lock()
			result[clan.Tag] = hexColor
			mu.Unlock()
		}

		return nil
	})
	if err != nil {
		b.logger.Warn(err, nil)
	}

	return result
}

func (b *Battle) fetchClanLanguage(clanInfoArray []data.WGClansInfoData) map[string]string {
	result := make(map[string]string)

	// URLを検出する正規表現パターン
	urlPattern := `https?://[^\s]+`
	re := regexp.MustCompile(urlPattern)

	options := whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Jpn: true,
			whatlanggo.Kor: true,
			whatlanggo.Cmn: true,
		},
	}

	var mu sync.Mutex
	err := doParallel(clanInfoArray, func(clan data.WGClansInfoData) error {
		// URLを空文字に
		description := re.ReplaceAllString(clan.Description, "")
		// 改行を空文字に
		description = strings.ReplaceAll(description, "\n", "")

		if len(description) == 0 {
			return nil
		}

		info := whatlanggo.DetectWithOptions(description, options)

		mu.Lock()
		result[clan.Tag] = info.Lang.Iso6391()
		mu.Unlock()

		return nil
	})
	if err != nil {
		b.logger.Warn(err, nil)
	}

	return result
}

func (b *Battle) compose(
	tempArenaInfo data.TempArenaInfo,
	accountInfo data.WGAccountInfo,
	accountList data.WGAccountList,
	clans data.Clans,
	allPlayerShipsStats data.AllPlayerShipsStats,
	warships data.Warships,
	allExpectedStats data.ExpectedStats,
	battleArenas data.WGBattleArenas,
	battleTypes data.WGBattleTypes,
) data.Battle {
	friends := make(data.Players, 0)
	enemies := make(data.Players, 0)

	var ownShip string

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clans[accountID]

		warship, ok := warships[vehicle.ShipID]
		if !ok {
			warship = data.Warship{
				Name:   "Unknown",
				Tier:   0,
				Type:   data.ShipTypeNONE,
				Nation: "",
			}
		}
		if nickname == tempArenaInfo.PlayerName {
			ownShip = warship.Name
		}

		stats := data.NewStats(
			vehicle.ShipID,
			accountInfo[accountID],
			allPlayerShipsStats.Player(accountID),
			allExpectedStats,
			warships,
			tempArenaInfo,
		)

		player := data.Player{
			PlayerInfo: data.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo[accountID].HiddenProfile,
			},
			ShipInfo: data.ShipInfo{
				ID:        vehicle.ShipID,
				Name:      warship.Name,
				Nation:    warship.Nation,
				Tier:      warship.Tier,
				Type:      warship.Type,
				IsPremium: warship.IsPremium,
				AvgDamage: allExpectedStats[vehicle.ShipID].AverageDamageDealt,
			},
			PvPSolo:  playerStats(data.StatsPatternPvPSolo, stats, accountID, vehicle.ShipID, tempArenaInfo, warships),
			PvPAll:   playerStats(data.StatsPatternPvPAll, stats, accountID, vehicle.ShipID, tempArenaInfo, warships),
			RankSolo: playerStats(data.StatsPatternRankSolo, stats, accountID, vehicle.ShipID, tempArenaInfo, warships),
		}

		if vehicle.IsFriend() {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Sort(friends)
	sort.Sort(enemies)

	teams := []data.Team{
		{Players: friends},
		{Players: enemies},
	}

	battle := data.Battle{
		Meta: data.Meta{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    tempArenaInfo.BattleArena(battleArenas),
			Type:     tempArenaInfo.BattleType(battleTypes),
			OwnShip:  ownShip,
		},
		Teams: teams,
	}

	return battle
}

func playerStats(
	statsPattern data.StatsPattern,
	stats *data.Stats,
	accountID int,
	shipID int,
	tempArenaInfo data.TempArenaInfo,
	warships data.Warships,
) data.PlayerStats {
	threatLevel := yamibuka.CalculateThreatLevel(yamibuka.NewThreatLevelFactor(
		accountID,
		tempArenaInfo,
		warships,
		shipID,
		stats.Battles(data.StatsCategoryShip, statsPattern),
		stats.AvgDamage(data.StatsCategoryShip, statsPattern),
		stats.WinRate(data.StatsCategoryShip, statsPattern),
		stats.SurvivedRate(data.StatsCategoryShip, statsPattern),
		stats.PlanesKilled(data.StatsCategoryShip),
		stats.Battles(data.StatsCategoryOverall, statsPattern),
		stats.AvgDamage(data.StatsCategoryOverall, statsPattern),
		stats.WinRate(data.StatsCategoryOverall, statsPattern),
		stats.AvgKill(data.StatsCategoryOverall, statsPattern),
		stats.KdRate(data.StatsCategoryOverall, statsPattern),
	))

	return data.PlayerStats{
		ShipStats: data.ShipStats{
			Battles:            stats.Battles(data.StatsCategoryShip, statsPattern),
			Damage:             stats.AvgDamage(data.StatsCategoryShip, statsPattern),
			MaxDamage:          stats.MaxDamage(data.StatsCategoryShip, statsPattern),
			WinRate:            stats.WinRate(data.StatsCategoryShip, statsPattern),
			WinSurvivedRate:    stats.WinSurvivedRate(data.StatsCategoryShip, statsPattern),
			LoseSurvivedRate:   stats.LoseSurvivedRate(data.StatsCategoryShip, statsPattern),
			KdRate:             stats.KdRate(data.StatsCategoryShip, statsPattern),
			Kill:               stats.AvgKill(data.StatsCategoryShip, statsPattern),
			Exp:                stats.AvgExp(data.StatsCategoryShip, statsPattern),
			PR:                 stats.PR(data.StatsCategoryShip, statsPattern),
			MainBatteryHitRate: stats.MainBatteryHitRate(statsPattern),
			TorpedoesHitRate:   stats.TorpedoesHitRate(statsPattern),
			PlanesKilled:       stats.PlanesKilled(statsPattern),
			PlatoonRate:        stats.PlatoonRate(data.StatsCategoryShip),
		},
		OverallStats: data.OverallStats{
			Battles:           stats.Battles(data.StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(data.StatsCategoryOverall, statsPattern),
			MaxDamage:         stats.MaxDamage(data.StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(data.StatsCategoryOverall, statsPattern),
			WinSurvivedRate:   stats.WinSurvivedRate(data.StatsCategoryOverall, statsPattern),
			LoseSurvivedRate:  stats.LoseSurvivedRate(data.StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(data.StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(data.StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(data.StatsCategoryOverall, statsPattern),
			PR:                stats.PR(data.StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern),
			UsingTierRate:     stats.UsingTierRate(statsPattern),
			PlatoonRate:       stats.PlatoonRate(data.StatsCategoryOverall),
			ThreatLevel:       threatLevel,
		},
	}
}
