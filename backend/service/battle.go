package service

import (
	"sort"
	"sync"

	"changeme/backend/domain"
	"changeme/backend/infra"
	"changeme/backend/vo"
)

type Battle struct {
	parallels         uint
	wargaming         infra.WargamingInterface
	numbers           infra.NumbersInterface
	unregistered      infra.UnregisteredInterface
	tempArenaInfoRepo infra.TempArenaInfoInterface
	isFirstBattle     bool

	warship       map[int]vo.Warship
	expectedStats vo.NSExpectedStats
	battleArenas  vo.WGBattleArenas
	battleTypes   vo.WGBattleTypes
}

func NewBattle(
	parallels uint,
	wargaming infra.WargamingInterface,
	tempArenaInfoRepo infra.TempArenaInfoInterface,
	numbers infra.NumbersInterface,
	unregistered infra.UnregisteredInterface,
) *Battle {
	return &Battle{
		parallels:         parallels,
		wargaming:         wargaming,
		tempArenaInfoRepo: tempArenaInfoRepo,
		numbers:           numbers,
		unregistered:      unregistered,
		isFirstBattle:     true,
	}
}

func (b *Battle) Battle(userConfig vo.UserConfig) (vo.Battle, error) {
	b.wargaming.SetAppID(userConfig.Appid)
	var result vo.Battle

	// Fetch on-memory stored data
	warshipResult := make(chan vo.Result[map[int]vo.Warship])
	expectedStatsResult := make(chan vo.Result[vo.NSExpectedStats])
	battleArenasResult := make(chan vo.Result[vo.WGBattleArenas])
	battleTypesResult := make(chan vo.Result[vo.WGBattleTypes])
	if b.isFirstBattle {
		go b.fetchWarship(warshipResult)
		go b.fetchExpectedStats(expectedStatsResult)
		go b.fetchBattleArenas(battleArenasResult)
		go b.fetchBattleTypes(battleTypesResult)
	}

	// Get tempArenaInfo.json
	tempArenaInfo, err := b.tempArenaInfoRepo.Get(userConfig.InstallPath)
	if err != nil {
		return result, err
	}
	if userConfig.SaveTempArenaInfo {
		if err := b.tempArenaInfoRepo.Save(tempArenaInfo); err != nil {
			return result, err
		}
	}

	// Get Account ID list
	accountNames := tempArenaInfo.AccountNames()
	accountList, err := b.wargaming.AccountList(accountNames)
	if err != nil {
		return result, err
	}
	accountIDs := accountList.AccountIDs()

	// Fetch each stats
	accountInfoResult := make(chan vo.Result[vo.WGAccountInfo])
	shipStatsResult := make(chan vo.Result[map[int]vo.WGShipsStats])
	clanResult := make(chan vo.Result[map[int]vo.Clan])
	go b.accountInfo(accountIDs, accountInfoResult)
	go b.shipStats(accountIDs, shipStatsResult)
	go b.clanTag(accountIDs, clanResult)

	errs := make([]error, 0)

	if b.isFirstBattle {
		warship := <-warshipResult
		b.warship = warship.Value
		errs = append(errs, warship.Error)

		expectedStats := <-expectedStatsResult
		b.expectedStats = expectedStats.Value
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
		b.expectedStats,
		b.battleArenas,
		b.battleTypes,
	)

	b.isFirstBattle = false

	return result, nil
}

func (b *Battle) fetchWarship(result chan vo.Result[map[int]vo.Warship]) {
	warships := make(map[int]vo.Warship, 0)

	res, err := b.wargaming.EncycShips(1)
	if err != nil {
		result <- vo.Result[map[int]vo.Warship]{Value: warships, Error: err}
		return
	}

	var mu sync.Mutex
	pages := makeRange(1, res.Meta.PageTotal)
	err = doParallel(b.parallels, pages, func(page int) error {
		encyclopediaShips, err := b.wargaming.EncycShips(page)
		if err != nil {
			return err
		}

		for shipID, warship := range encyclopediaShips.Data {
			mu.Lock()
			warships[shipID] = vo.Warship{
				Name:   warship.Name,
				Tier:   warship.Tier,
				Type:   vo.NewShipType(warship.Type),
				Nation: warship.Nation,
			}
			mu.Unlock()
		}

		return nil
	})
	if err != nil {
		result <- vo.Result[map[int]vo.Warship]{Value: warships, Error: err}
		return
	}

	unregisteredShipInfo, err := b.unregistered.Warship()
	if err != nil {
		result <- vo.Result[map[int]vo.Warship]{Value: warships, Error: err}
		return
	}
	for k, v := range unregisteredShipInfo {
		if _, ok := warships[k]; !ok {
			warships[k] = v
		}
	}

	result <- vo.Result[map[int]vo.Warship]{Value: warships, Error: nil}
}

func (b *Battle) fetchExpectedStats(result chan vo.Result[vo.NSExpectedStats]) {
	expectedStats, err := b.numbers.ExpectedStats()
	result <- vo.Result[vo.NSExpectedStats]{Value: expectedStats, Error: err}
}

func (b *Battle) fetchBattleArenas(result chan vo.Result[vo.WGBattleArenas]) {
	battleArenas, err := b.wargaming.BattleArenas()
	result <- vo.Result[vo.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (b *Battle) fetchBattleTypes(result chan vo.Result[vo.WGBattleTypes]) {
	battleTypes, err := b.wargaming.BattleTypes()
	result <- vo.Result[vo.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (b *Battle) accountInfo(accountIDs []int, result chan vo.Result[vo.WGAccountInfo]) {
	accountInfo, err := b.wargaming.AccountInfo(accountIDs)
	result <- vo.Result[vo.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (b *Battle) shipStats(accountIDs []int, result chan vo.Result[map[int]vo.WGShipsStats]) {
	shipStatsMap := make(map[int]vo.WGShipsStats)
	var mu sync.Mutex
	err := doParallel(b.parallels, accountIDs, func(accountID int) error {
		shipStats, err := b.wargaming.ShipsStats(accountID)
		if err != nil {
			return err
		}

		mu.Lock()
		shipStatsMap[accountID] = shipStats
		mu.Unlock()

		return nil
	})

	result <- vo.Result[map[int]vo.WGShipsStats]{Value: shipStatsMap, Error: err}
}

func (b *Battle) clanTag(accountIDs []int, result chan vo.Result[map[int]vo.Clan]) {
	clanMap := make(map[int]vo.Clan)

	clansAccountInfo, err := b.wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[map[int]vo.Clan]{Value: clanMap, Error: err}

		return
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := b.wargaming.ClansInfo(clanIDs)
	if err != nil {
		result <- vo.Result[map[int]vo.Clan]{Value: clanMap, Error: err}

		return
	}

	for _, accountID := range accountIDs {
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		clanMap[accountID] = vo.Clan{Tag: clanTag, ID: clanID}
	}

	result <- vo.Result[map[int]vo.Clan]{Value: clanMap, Error: nil}
}

func (b *Battle) compose(
	tempArenaInfo vo.TempArenaInfo,
	accountInfo vo.WGAccountInfo,
	accountList vo.WGAccountList,
	clan map[int]vo.Clan,
	shipStats map[int]vo.WGShipsStats,
	warships map[int]vo.Warship,
	expectedStats vo.NSExpectedStats,
	battleArenas vo.WGBattleArenas,
	battleTypes vo.WGBattleTypes,
) vo.Battle {
	friends := make(vo.Players, 0)
	enemies := make(vo.Players, 0)

	var ownShip string

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clan[accountID]

		warship, ok := warships[vehicle.ShipID]
		if !ok {
			warship = vo.Warship{
				Name:   "Unknown",
				Tier:   0,
				Type:   vo.NONE,
				Nation: "",
			}
		}
		if nickname == tempArenaInfo.PlayerName {
			ownShip = warship.Name
		}

		stats := domain.Stats{
			AccountInfo: accountInfo.Data[accountID],
			Expected:    expectedStats.Data[vehicle.ShipID],
		}
		for _, v := range shipStats[accountID].Data[accountID] {
			if v.ShipID == vehicle.ShipID {
				stats.SetShipStats(v)
				break
			}
		}

		player := vo.Player{
			ShipInfo: vo.ShipInfo{
				ID:        vehicle.ShipID,
				Name:      warship.Name,
				Nation:    warship.Nation,
				Tier:      warship.Tier,
				Type:      warship.Type,
				AvgDamage: stats.Expected.AverageDamageDealt,
			},
			ShipStats: vo.ShipStats{
				Battles:            stats.Battles(domain.ModeShip),
				Damage:             stats.AvgDamage(domain.ModeShip),
				WinRate:            stats.WinRate(domain.ModeShip),
				WinSurvivedRate:    stats.WinSurvivedRate(domain.ModeShip),
				LoseSurvivedRate:   stats.LoseSurvivedRate(domain.ModeShip),
				KdRate:             stats.KdRate(domain.ModeShip),
				Exp:                stats.AvgExp(domain.ModeShip),
				MainBatteryHitRate: stats.MainBatteryHitRate(domain.ModeShip),
				TorpedoesHitRate:   stats.TorpedoesHitRate(domain.ModeShip),
				PR:                 stats.ShipPR(),
			},
			PlayerInfo: vo.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: stats.AccountInfo.HiddenProfile,
			},
			OverallStats: vo.OverallStats{
				Battles:           stats.Battles(domain.ModeOverall),
				Damage:            stats.AvgDamage(domain.ModeOverall),
				WinRate:           stats.WinRate(domain.ModeOverall),
				WinSurvivedRate:   stats.WinSurvivedRate(domain.ModeOverall),
				LoseSurvivedRate:  stats.LoseSurvivedRate(domain.ModeOverall),
				KdRate:            stats.KdRate(domain.ModeOverall),
				Exp:               stats.AvgExp(domain.ModeOverall),
				AvgTier:           stats.AvgTier(accountID, warships, shipStats),
				UsingShipTypeRate: stats.UsingShipTypeRate(accountID, warships, shipStats),
				UsingTierRate:     stats.UsingTierRate(accountID, warships, shipStats),
			},
		}

		if vehicle.Relation == 0 || vehicle.Relation == 1 {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Sort(friends)
	sort.Sort(enemies)

	teams := []vo.Team{
		{
			Players: friends,
			Name:    "味方チーム",
		},
		{
			Players: enemies,
			Name:    "敵チーム",
		},
	}

	battle := vo.Battle{
		Meta: vo.Meta{
			Date:    tempArenaInfo.FormattedDateTime(),
			Arena:   tempArenaInfo.BattleArena(battleArenas),
			Type:    tempArenaInfo.BattleType(battleTypes),
			OwnShip: ownShip,
		},
		Teams: teams,
	}

	return battle
}
