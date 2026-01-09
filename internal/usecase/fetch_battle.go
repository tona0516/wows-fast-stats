package usecase

import (
	"context"
	"sort"
	"wfs/internal/adapter"
	"wfs/internal/data"
	"wfs/internal/service"

	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

type FetchBattle struct {
	userDataFetcher    *service.UserDataFetcher
	nonUserDataFetcher *service.NonUserDataFetcher
	wails              adapter.Wails
	cacheStore         adapter.CacheStore
	logger             adapter.Logger
}

func NewFetchBattle(i do.Injector) (*FetchBattle, error) {
	return &FetchBattle{
		userDataFetcher:    do.MustInvoke[*service.UserDataFetcher](i),
		nonUserDataFetcher: do.MustInvoke[*service.NonUserDataFetcher](i),
		wails:              do.MustInvoke[adapter.Wails](i),
		cacheStore:         do.MustInvoke[adapter.CacheStore](i),
		logger:             do.MustInvoke[adapter.Logger](i),
	}, nil
}

func (b *FetchBattle) Invoke(
	ctx context.Context,
	tempArenaInfo data.TempArenaInfo,
) {
	_ = b.cacheStore.SetOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	eg := errgroup.Group{}

	var userData *data.UserData
	eg.Go(func() error {
		var err error
		userData, err = b.userDataFetcher.Fetch(tempArenaInfo.AccountNames())
		return err
	})

	var nonUserData *data.NonUserData
	eg.Go(func() error {
		var err error
		nonUserData, err = b.nonUserDataFetcher.Fetch()
		return err
	})

	if err := eg.Wait(); err != nil {
		b.wails.EmitEvent(ctx, EventErr, err)
		return
	}

	result := b.compose(
		tempArenaInfo,
		userData,
		nonUserData,
	)

	b.wails.EmitEvent(ctx, EventFetchDone, result)
}

func (b *FetchBattle) compose(
	tempArenaInfo data.TempArenaInfo,
	userData *data.UserData,
	nonUserData *data.NonUserData,
) data.Battle {
	friends := make(data.Players, 0)
	enemies := make(data.Players, 0)

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := userData.AccountList.AccountID(nickname)
		clan := userData.Clans[accountID]

		warship, ok := nonUserData.Warships[vehicle.ShipID]
		if !ok {
			warship = *data.NewUnknownWarship()
		}

		stats := data.NewPersonalStats(
			vehicle.ShipID,
			userData.AccountInfo.Data[accountID],
			userData.AllPlayerShipsStats.Player(accountID),
			userData.AllPlayerShipsBadges[accountID],
			nonUserData.Warships,
			tempArenaInfo,
		)

		player := data.Player{
			PlayerInfo: data.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: userData.AccountInfo.Data[accountID].HiddenProfile,
			},
			Warship: warship,
			PvPSolo: data.BuildPlayerStats(
				data.StatsPatternPvPSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				nonUserData.Warships,
			),
			PvPAll: data.BuildPlayerStats(
				data.StatsPatternPvPAll,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				nonUserData.Warships,
			),
			RankSolo: data.BuildPlayerStats(
				data.StatsPatternRankSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				nonUserData.Warships,
			),
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
				TeamAverageStats: data.CalculateTeamAverageStats(friends, data.StatsPatternPvPAll),
				TeamThreatLevel:  data.CalculateTeamThreatLevel(friends, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamAverageStats: data.CalculateTeamAverageStats(friends, data.StatsPatternPvPSolo),
				TeamThreatLevel:  data.CalculateTeamThreatLevel(friends, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamAverageStats: data.CalculateTeamAverageStats(friends, data.StatsPatternRankSolo),
				TeamThreatLevel:  data.CalculateTeamThreatLevel(friends, data.StatsPatternRankSolo),
			},
		},
		{
			Players: enemies,
			PvPAll: data.TeamStats{
				TeamAverageStats: data.CalculateTeamAverageStats(enemies, data.StatsPatternPvPAll),
				TeamThreatLevel:  data.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamAverageStats: data.CalculateTeamAverageStats(enemies, data.StatsPatternPvPSolo),
				TeamThreatLevel:  data.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamAverageStats: data.CalculateTeamAverageStats(enemies, data.StatsPatternRankSolo),
				TeamThreatLevel:  data.CalculateTeamThreatLevel(enemies, data.StatsPatternRankSolo),
			},
		},
	}

	battle := data.Battle{
		Meta: data.BattleMetaData{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    tempArenaInfo.BattleArena(nonUserData.BattleArenas),
			Type:     tempArenaInfo.BattleType(nonUserData.BattleTypes),
		},
		Teams: teams,
	}

	return battle
}
