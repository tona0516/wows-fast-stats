package usecase

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"sync"
	"wfs/internal/apperr"
	"wfs/internal/data"
	"wfs/internal/infra"
	"wfs/internal/service"
	"wfs/internal/util"
	"wfs/internal/yamibuka"

	"github.com/abadojack/whatlanggo"
	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

type FetchBattle struct {
	nonUserDataFetcher *service.NonUserDataFetcher
	localStorage       infra.LocalStorage
	wargaming          infra.WargamingApiClient
	uwargaming         infra.ClanApiClient
	numbers            infra.NumbersApiClient
	logger             infra.Logger
	eventsEmitFunc     eventEmitFunc
}

func NewFetchBattle(
	nonUserDataFetcher *service.NonUserDataFetcher,
	localStorage infra.LocalStorage,
	wargaming infra.WargamingApiClient,
	uwargaming infra.ClanApiClient,
	numbers infra.NumbersApiClient,
	logger infra.Logger,
	eventsEmitFunc eventEmitFunc,
) *FetchBattle {
	return &FetchBattle{
		nonUserDataFetcher: nonUserDataFetcher,
		localStorage:       localStorage,
		wargaming:          wargaming,
		uwargaming:         uwargaming,
		numbers:            numbers,
		logger:             logger,
		eventsEmitFunc:     eventsEmitFunc,
	}
}

func (b *FetchBattle) Invoke(
	ctx context.Context,
	tempArenaInfo data.TempArenaInfo,
) {
	eg := errgroup.Group{}

	var nonUserData *service.NonUserData
	eg.Go(func() error {
		var err error
		nonUserData, err = b.nonUserDataFetcher.Fetch()
		return err
	})

	_ = b.localStorage.SetOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	accountList, err := b.wargaming.AccountList(tempArenaInfo.AccountNames())
	if err != nil {
		b.eventsEmitFunc(ctx, EventErr, apperr.ToStringCode(err))
		return
	}
	accountIDs := accountList.AccountIDs()

	var accountInfo data.WGAccountInfo
	eg.Go(func() error {
		var err error
		accountInfo, err = b.wargaming.AccountInfo(accountIDs)
		return err
	})

	var allShipStats data.AllPlayerShipsStats
	eg.Go(func() error {
		var err error
		allShipStats, err = b.fetchAllPlayerShipsStats(accountIDs)
		return err
	})

	var allShipBadges data.AllPlayerShipsBadges
	eg.Go(func() error {
		var err error
		allShipBadges, err = b.fetchAllPlayerShipsBadges(accountIDs)
		return err
	})

	var clans data.Clans
	eg.Go(func() error {
		var err error
		clans, err = b.fetchClan(accountIDs)
		return err
	})

	if err := eg.Wait(); err != nil {
		b.eventsEmitFunc(ctx, EventErr, err)
		return
	}

	result := b.compose(
		tempArenaInfo,
		accountInfo,
		accountList,
		clans,
		allShipStats,
		allShipBadges,
		nonUserData.Warships,
		nonUserData.BattleArenas,
		nonUserData.BattleTypes,
	)

	b.eventsEmitFunc(ctx, EventFetchDone, result)
}

func (b *FetchBattle) fetchAllPlayerShipsStats(accountIDs []int) (data.AllPlayerShipsStats, error) {
	result := make(data.AllPlayerShipsStats)
	var mu sync.Mutex

	eg := errgroup.Group{}

	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := b.wargaming.ShipsStats(accountID)
			if err != nil {
				return failure.Wrap(err)
			}

			mu.Lock()
			result[accountID] = resp
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, failure.Wrap(err)
	}

	return result, nil
}

func (b *FetchBattle) fetchClan(accountIDs []int) (data.Clans, error) {
	result := make(data.Clans)

	clansAccountInfo, err := b.wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := b.wargaming.ClansInfo(clanIDs)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	clanInfoArray := clansInfo.ToArray()
	colorMap, err := b.fetchClanColor(clanInfoArray)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	languageMap := b.fetchClanLanguage(clanInfoArray)

	for _, accountID := range accountIDs {
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		hexColor := colorMap[clanTag]
		language := languageMap[clanTag]

		result[accountID] = data.Clan{
			ID:       clanID,
			Tag:      clanTag,
			HexColor: hexColor,
			Language: language,
		}
	}

	return result, nil
}

func (b *FetchBattle) fetchClanColor(clanInfoArray []data.WGClansInfoData) (map[string]string, error) {
	result := make(map[string]string)
	var mu sync.Mutex

	eg := errgroup.Group{}
	for _, clanInfo := range clanInfoArray {
		eg.Go(func() error {
			autocomplete, err := b.uwargaming.ClanAutoComplete(clanInfo.Tag)
			if err != nil {
				return err
			}

			hexColor := autocomplete.HexColor(clanInfo.Tag)
			if hexColor == "" {
				return nil
			}

			mu.Lock()
			result[clanInfo.Tag] = hexColor
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, failure.Wrap(err)
	}

	return result, nil
}

func (b *FetchBattle) fetchClanLanguage(clanInfoArray []data.WGClansInfoData) map[string]string {
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
	err := util.DoParallel(clanInfoArray, func(clan data.WGClansInfoData) error {
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
		b.logger.Error(err, nil)
	}

	return result
}

func (b *FetchBattle) fetchAllPlayerShipsBadges(accountIDs []int) (data.AllPlayerShipsBadges, error) {
	result := make(data.AllPlayerShipsBadges)
	var mu sync.Mutex

	eg := errgroup.Group{}
	for _, accountID := range accountIDs {
		eg.Go(func() error {
			shipsBadges, err := b.wargaming.ShipsBadges(accountID)
			if err != nil {
				return failure.Wrap(err)
			}

			mu.Lock()
			result[accountID] = shipsBadges.Data[accountID]
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, failure.Wrap(err)
	}

	return result, nil
}

func (b *FetchBattle) compose(
	tempArenaInfo data.TempArenaInfo,
	accountInfo data.WGAccountInfo,
	accountList data.WGAccountList,
	clans data.Clans,
	allPlayerShipsStats data.AllPlayerShipsStats,
	allPlayerShipsBadges data.AllPlayerShipsBadges,
	warships data.Warships,
	battleArenas map[int]string,
	battleTypes map[string]string,
) data.Battle {
	friends := make(data.Players, 0)
	enemies := make(data.Players, 0)

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clans[accountID]

		warship, ok := warships[vehicle.ShipID]
		if !ok {
			println("unknown ship: ", vehicle.ShipID)
			warship = *data.NewUnknownWarship()
		}

		stats := data.NewPersonalStats(
			vehicle.ShipID,
			accountInfo.Data[accountID],
			allPlayerShipsStats.Player(accountID),
			allPlayerShipsBadges[accountID],
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
			Warship:  warship,
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
				TeamAverageStats: calculateTeamAverageStats(friends, data.StatsPatternPvPAll),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(friends, data.StatsPatternPvPSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(friends, data.StatsPatternRankSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternRankSolo),
			},
		},
		{
			Players: enemies,
			PvPAll: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(enemies, data.StatsPatternPvPAll),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(enemies, data.StatsPatternPvPSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(enemies, data.StatsPatternRankSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternRankSolo),
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

func calculateTeamAverageStats(
	players data.Players,
	statsPattern data.StatsPattern,
) data.TeamAverageStats {
	var shipPRSum, shipDamageSum, shipWinRateSum float64
	var shipBattlesSum, shipStatsCount uint

	var overallPRSum, overallDamageSum, overallWinRateSum float64
	var overallBattlesSum, overallStatsCount uint

	for _, player := range players {
		var shipStats data.ShipStats
		var overallStats data.OverallStats
		switch statsPattern {
		case data.StatsPatternPvPSolo:
			shipStats = player.PvPSolo.ShipStats
			overallStats = player.PvPSolo.OverallStats
		case data.StatsPatternPvPAll:
			shipStats = player.PvPAll.ShipStats
			overallStats = player.PvPAll.OverallStats
		case data.StatsPatternRankSolo:
			shipStats = player.RankSolo.ShipStats
			overallStats = player.RankSolo.OverallStats
		}

		if shipStats.Battles > 0 {
			shipPRSum += shipStats.PR.Value
			shipDamageSum += shipStats.Damage.Value
			shipWinRateSum += shipStats.WinRate.Value
			shipBattlesSum += shipStats.Battles

			shipStatsCount++
		}

		if overallStats.Battles > 0 {
			overallPRSum += overallStats.PR.Value
			overallDamageSum += overallStats.Damage.Value
			overallWinRateSum += overallStats.WinRate.Value
			overallBattlesSum += overallStats.Battles

			overallStatsCount++
		}
	}

	return data.TeamAverageStats{
		ShipPR:         util.SafeDivide(shipPRSum, shipStatsCount),
		ShipDamage:     util.SafeDivide(shipDamageSum, shipStatsCount),
		ShipWinRate:    util.SafeDivide(shipWinRateSum, shipStatsCount),
		ShipBattles:    uint(util.SafeDivide(float64(shipBattlesSum), shipStatsCount)),
		OverallPR:      util.SafeDivide(overallPRSum, overallStatsCount),
		OverallDamage:  util.SafeDivide(overallDamageSum, overallStatsCount),
		OverallWinRate: util.SafeDivide(overallWinRateSum, overallStatsCount),
		OverallBattles: uint(util.SafeDivide(float64(overallBattlesSum), overallStatsCount)),
	}
}
