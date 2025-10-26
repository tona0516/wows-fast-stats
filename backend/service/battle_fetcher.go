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

type BattleFetcher struct {
	ctx            context.Context
	wargaming      repository.WargamingInterface
	uwargaming     repository.UnofficialWargamingInterface
	numbers        repository.NumbersInterface
	persistence    repository.PersistenceInterface
	logger         repository.LoggerInterface
	eventsEmitFunc eventEmitFunc

	isFirstBattle                      bool
	isNotifyExpectedStatsUnavaillalble bool
	warship                            data.Warships
	allExpectedStats                   data.ExpectedStats
	battleArenas                       data.WGBattleArenas
	battleTypes                        data.WGBattleTypes
}

func NewBattleFetcher(
	ctx context.Context,
	wargaming repository.WargamingInterface,
	uwargaming repository.UnofficialWargamingInterface,
	numbers repository.NumbersInterface,
	persistence repository.PersistenceInterface,
	logger repository.LoggerInterface,
	eventsEmitFunc eventEmitFunc,
) *BattleFetcher {
	return &BattleFetcher{
		ctx:                                ctx,
		wargaming:                          wargaming,
		uwargaming:                         uwargaming,
		numbers:                            numbers,
		persistence:                        persistence,
		logger:                             logger,
		eventsEmitFunc:                     eventsEmitFunc,
		isFirstBattle:                      true,
		isNotifyExpectedStatsUnavaillalble: false,
	}
}

func (b *BattleFetcher) Invoke(tempArenaInfo data.TempArenaInfo) {
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

	// persist own ign for reporting
	_ = b.persistence.SaveOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	accountList, err := b.wargaming.AccountList(tempArenaInfo.AccountNames())
	if err != nil {
		b.eventsEmitFunc(b.ctx, EventErr, apperr.ToStringCode(err))
		return
	}
	accountIDs := accountList.AccountIDs()

	// Fetch each stats
	accountInfoResult := make(chan data.Result[data.WGAccountInfo])
	shipStatsResult := make(chan data.Result[data.AllPlayerShipsStats])
	clanResult := make(chan data.Result[data.Clans])
	shipsBadgesResult := make(chan data.Result[data.AllPlayerShipsBadges])
	go b.fetchAccountInfo(accountIDs, accountInfoResult)
	go b.fetchAllPlayerShipsStats(accountIDs, shipStatsResult)
	go b.fetchAllPlayerShipsBadges(accountIDs, shipsBadgesResult)
	go b.fetchClan(accountIDs, clanResult)

	errs := make([]error, 0)
	if b.isFirstBattle {
		b.eventsEmitFunc(b.ctx, EventFetchOthers, nil)

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

	b.eventsEmitFunc(b.ctx, EventFetchPlayers, nil)

	accountInfo := <-accountInfoResult
	errs = append(errs, accountInfo.Error)

	shipStats := <-shipStatsResult
	errs = append(errs, shipStats.Error)

	shipsBadges := <-shipsBadgesResult
	errs = append(errs, shipsBadges.Error)

	clan := <-clanResult
	errs = append(errs, clan.Error)

	for _, err := range errs {
		if err != nil {
			if failure.Is(err, apperr.ExpectedStatsUnavaillalble) && !b.isNotifyExpectedStatsUnavaillalble {
				b.eventsEmitFunc(b.ctx, EventErr, apperr.ExpectedStatsUnavaillalble.ErrorCode())
				b.isNotifyExpectedStatsUnavaillalble = true
				continue
			}

			b.eventsEmitFunc(b.ctx, EventErr, apperr.ToStringCode(err))
			return
		}
	}

	result = b.compose(
		tempArenaInfo,
		accountInfo.Value,
		accountList,
		clan.Value,
		shipStats.Value,
		shipsBadges.Value,
		b.warship,
		b.allExpectedStats,
		b.battleArenas,
		b.battleTypes,
	)

	b.isFirstBattle = false
	b.eventsEmitFunc(b.ctx, EventFetchDone, result)
}

func (b *BattleFetcher) fetchWarships(channel chan data.Result[data.Warships]) {
	warships := make(data.Warships)
	var result data.Result[data.Warships]

	var mu sync.Mutex

	fetch := func(page int) (int, error) {
		res, err := b.wargaming.EncycShips(page)
		if err != nil {
			return 0, err
		}

		for shipID, warship := range res.Data {
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
		return res.Meta.PageTotal, nil
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

	result.Value = warships
	channel <- result
}

func (b *BattleFetcher) fetchExpectedStats(channel chan data.Result[data.ExpectedStats]) {
	var result data.Result[data.ExpectedStats]

	// 最新の予測成績を取得
	expectedStats, errFetch := b.numbers.ExpectedStats()
	if errFetch == nil {
		_ = b.persistence.SaveExpectedStats(expectedStats)

		result.Value = expectedStats
		channel <- result
		return
	}

	// 取得できない場合、キャッシュを利用する
	expectedStatsCache, errCache := b.persistence.LoadExpectedStats()
	if errCache == nil {
		result.Value = expectedStatsCache
		channel <- result
		return
	}

	// APIからもキャッシュも取得できない場合、エラーを返す
	result.Error = failure.New(apperr.ExpectedStatsUnavaillalble, failure.Context{
		"err_fetch": errFetch.Error(),
		"err_cache": errCache.Error(),
	})
	channel <- result
}

func (b *BattleFetcher) fetchBattleArenas(channel chan data.Result[data.WGBattleArenas]) {
	battleArenas, err := b.wargaming.BattleArenas()
	channel <- data.Result[data.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (b *BattleFetcher) fetchBattleTypes(channel chan data.Result[data.WGBattleTypes]) {
	battleTypes, err := b.wargaming.BattleTypes()
	channel <- data.Result[data.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (b *BattleFetcher) fetchAccountInfo(accountIDs []int, channel chan data.Result[data.WGAccountInfo]) {
	accountInfo, err := b.wargaming.AccountInfo(accountIDs)
	channel <- data.Result[data.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (b *BattleFetcher) fetchAllPlayerShipsStats(
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

func (b *BattleFetcher) fetchClan(accountIDs []int, channel chan data.Result[data.Clans]) {
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
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		hexColor := colorMap[clanTag]
		language := languageMap[clanTag]

		clans[accountID] = data.Clan{Tag: clanTag, ID: clanID, HexColor: hexColor, Language: language}
	}

	result.Value = clans
	channel <- result
}

func (b *BattleFetcher) fetchClanColor(clanInfoArray []data.WGClansInfoData) map[string]string {
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

func (b *BattleFetcher) fetchClanLanguage(clanInfoArray []data.WGClansInfoData) map[string]string {
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

func (b *BattleFetcher) fetchAllPlayerShipsBadges(
	accountIDs []int,
	channel chan data.Result[data.AllPlayerShipsBadges],
) {
	shipsBadgesMap := make(data.AllPlayerShipsBadges)
	var mu sync.Mutex
	err := doParallel(accountIDs, func(accountID int) error {
		shipsBadges, err := b.wargaming.ShipsBadges(accountID)
		if err != nil {
			return err
		}

		mu.Lock()
		shipsBadgesMap[accountID] = shipsBadges.Data[accountID]
		mu.Unlock()

		return nil
	})

	channel <- data.Result[data.AllPlayerShipsBadges]{Value: shipsBadgesMap, Error: err}
}

func (b *BattleFetcher) compose(
	tempArenaInfo data.TempArenaInfo,
	accountInfo data.WGAccountInfo,
	accountList data.WGAccountList,
	clans data.Clans,
	allPlayerShipsStats data.AllPlayerShipsStats,
	allPlayerShipsBadges data.AllPlayerShipsBadges,
	warships data.Warships,
	allExpectedStats data.ExpectedStats,
	battleArenas data.WGBattleArenas,
	battleTypes data.WGBattleTypes,
) data.Battle {
	friends := make(data.Players, 0)
	enemies := make(data.Players, 0)

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clans[accountID]

		warship, ok := warships[vehicle.ShipID]
		if !ok {
			warship = data.Warship{
				Name:   "UNKNOWN",
				Tier:   0,
				Type:   data.ShipTypeNONE,
				Nation: "",
			}
		}

		stats := data.NewPersonalStats(
			vehicle.ShipID,
			accountInfo.Data[accountID],
			allPlayerShipsStats.Player(accountID),
			allPlayerShipsBadges[accountID],
			allExpectedStats,
			warships,
			tempArenaInfo,
		)

		player := data.Player{
			PlayerInfo: data.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo.Data[accountID].HiddenProfile,
			},
			ShipInfo: data.ShipInfo{
				ID:            vehicle.ShipID,
				Name:          warship.Name,
				Nation:        warship.Nation,
				Tier:          warship.Tier,
				Type:          warship.Type,
				IsPremium:     warship.IsPremium,
				AvgDamage:     allExpectedStats[vehicle.ShipID].AverageDamageDealt,
				DamageRatings: data.NewDamageRatings(allExpectedStats[vehicle.ShipID].AverageDamageDealt),
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
		{
			Players: friends,
			PvPAll: data.TeamStats{
				TeamThreatLevel: yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamThreatLevel: yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamThreatLevel: yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternRankSolo),
			},
		},
		{
			Players: enemies,
			PvPAll: data.TeamStats{
				TeamThreatLevel: yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamThreatLevel: yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamThreatLevel: yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternRankSolo),
			},
		},
	}

	battle := data.Battle{
		Meta: data.BattleMetaData{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    tempArenaInfo.BattleArena(battleArenas),
			Type:     tempArenaInfo.BattleType(battleTypes),
		},
		Teams: teams,
	}

	return battle
}

func playerStats(
	statsPattern data.StatsPattern,
	stats *data.PersonalStats,
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
		stats.AvgDamage(data.StatsCategoryShip, statsPattern).Value,
		stats.WinRate(data.StatsCategoryShip, statsPattern).Value,
		stats.SurvivedRate(data.StatsCategoryShip, statsPattern).All,
		stats.PlanesKilled(data.StatsCategoryShip),
		stats.Battles(data.StatsCategoryOverall, statsPattern),
		stats.AvgDamage(data.StatsCategoryOverall, statsPattern).Value,
		stats.WinRate(data.StatsCategoryOverall, statsPattern).Value,
		stats.AvgKill(data.StatsCategoryOverall, statsPattern),
		stats.KdRate(data.StatsCategoryOverall, statsPattern),
	))

	return data.PlayerStats{
		ShipStats: data.ShipStats{
			Battles:         stats.Battles(data.StatsCategoryShip, statsPattern),
			Damage:          stats.AvgDamage(data.StatsCategoryShip, statsPattern),
			MaxDamage:       stats.MaxDamage(data.StatsCategoryShip, statsPattern),
			WinRate:         stats.WinRate(data.StatsCategoryShip, statsPattern),
			SurvivedRate:    stats.SurvivedRate(data.StatsCategoryShip, statsPattern),
			KdRate:          stats.KdRate(data.StatsCategoryShip, statsPattern),
			Kill:            stats.AvgKill(data.StatsCategoryShip, statsPattern),
			Exp:             stats.AvgExp(data.StatsCategoryShip, statsPattern),
			PR:              stats.PR(data.StatsCategoryShip, statsPattern),
			HitRate:         stats.HitRate(statsPattern),
			PlanesKilled:    stats.PlanesKilled(statsPattern),
			PlatoonRate:     stats.PlatoonRate(data.StatsCategoryShip),
			EfficiencyBadge: stats.EfficiencyBadge(),
		},
		OverallStats: data.OverallStats{
			Battles:           stats.Battles(data.StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(data.StatsCategoryOverall, statsPattern),
			MaxDamage:         stats.MaxDamage(data.StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(data.StatsCategoryOverall, statsPattern),
			SurvivedRate:      stats.SurvivedRate(data.StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(data.StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(data.StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(data.StatsCategoryOverall, statsPattern),
			PR:                stats.PR(data.StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern),
			UsingTierRate:     stats.UsingTierRate(statsPattern),
			PlatoonRate:       stats.PlatoonRate(data.StatsCategoryOverall),
			EfficiencyBadge:   stats.EfficiencyBadges(),
			ThreatLevel:       threatLevel,
		},
	}
}
