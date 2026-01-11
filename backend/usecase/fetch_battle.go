package usecase

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"sync"
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/abadojack/whatlanggo"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

// URLを検出する正規表現パターン.
var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

type FetchBattle struct {
	wargamingClient adapter.WargamingClient
	clanClient      adapter.ClanClient
	wails           adapter.Wails
	cacheStore      adapter.CacheStore
	logger          adapter.Logger
}

func NewFetchBattle(i do.Injector) (*FetchBattle, error) {
	return &FetchBattle{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		clanClient:      do.MustInvoke[adapter.ClanClient](i),
		wails:           do.MustInvoke[adapter.Wails](i),
		cacheStore:      do.MustInvoke[adapter.CacheStore](i),
		logger:          do.MustInvoke[adapter.Logger](i),
	}, nil
}

func (b *FetchBattle) Invoke(
	ctx context.Context,
	tempArenaInfo data.TempArenaInfo,
	prefetchResult *data.PrefetchResult,
) {
	_ = b.cacheStore.SetOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	accountNames := tempArenaInfo.AccountNames()

	var accountList data.WGAccountList
	var err error
	measure("AccountList", func() {
		accountList, err = b.wargamingClient.AccountList(accountNames)
	})
	if err != nil {
		b.wails.EmitEvent(ctx, EventErr, err)
		return
	}
	accountIDs := accountList.AccountIDs()

	eg := errgroup.Group{}

	var accountInfo data.WGAccountInfo
	eg.Go(func() error {
		var err error
		measure("AccountInfo", func() {
			accountInfo, err = b.wargamingClient.AccountInfo(accountIDs)
		})
		return err
	})

	var allShipStats data.AllPlayerShipsStats
	eg.Go(func() error {
		var err error
		measure("fetchAllPlayerShipsStats", func() {
			allShipStats, err = b.fetchAllPlayerShipsStats(accountIDs)
		})
		return err
	})

	var allShipBadges data.AllPlayerShipsBadges
	eg.Go(func() error {
		var err error
		measure("fetchAllPlayerShipsBadges", func() {
			allShipBadges, err = b.fetchAllPlayerShipsBadges(accountIDs)
		})
		return err
	})

	var clans data.Clans
	eg.Go(func() error {
		var err error
		measure("fetchClan", func() {
			clans, err = b.fetchClan(accountIDs)
		})
		return err
	})

	if err := eg.Wait(); err != nil {
		b.wails.EmitEvent(ctx, EventErr, err)
		return
	}

	result := b.compose(
		prefetchResult,
		tempArenaInfo,
		accountInfo,
		accountList,
		clans,
		allShipStats,
		allShipBadges,
	)

	b.wails.EmitEvent(ctx, EventFetchDone, result)
}

func (b *FetchBattle) fetchAllPlayerShipsStats(accountIDs []int) (data.AllPlayerShipsStats, error) {
	result := make(data.AllPlayerShipsStats)
	var mu sync.Mutex

	eg := errgroup.Group{}

	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := b.wargamingClient.ShipsStats(accountID)
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

	clansAccountInfo, err := b.wargamingClient.ClansAccountInfo(accountIDs)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := b.wargamingClient.ClansInfo(clanIDs)
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

func (b *FetchBattle) fetchClanColor(clanInfoSlice []data.WGClansInfoData) (map[string]string, error) {
	result := make(map[string]string)
	var mu sync.Mutex

	eg := errgroup.Group{}
	for _, clanInfo := range clanInfoSlice {
		eg.Go(func() error {
			autocomplete, err := b.clanClient.ClanAutoComplete(clanInfo.Tag)
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

func (b *FetchBattle) fetchClanLanguage(clanInfoSilce []data.WGClansInfoData) map[string]string {
	result := make(map[string]string)

	options := whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Jpn: true,
			whatlanggo.Kor: true,
			whatlanggo.Cmn: true,
		},
	}

	for _, clanInfo := range clanInfoSilce {
		// URLを空文字に
		description := urlRegex.ReplaceAllString(clanInfo.Description, "")
		// 改行を空文字に
		description = strings.ReplaceAll(description, "\n", "")

		if len(description) == 0 {
			continue
		}

		info := whatlanggo.DetectWithOptions(description, options)
		result[clanInfo.Tag] = info.Lang.Iso6391()
	}

	return result
}

func (b *FetchBattle) fetchAllPlayerShipsBadges(accountIDs []int) (data.AllPlayerShipsBadges, error) {
	result := make(data.AllPlayerShipsBadges)
	var mu sync.Mutex

	eg := errgroup.Group{}
	for _, accountID := range accountIDs {
		eg.Go(func() error {
			shipsBadges, err := b.wargamingClient.ShipsBadges(accountID)
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
	prefetchResult *data.PrefetchResult,
	tempArenaInfo data.TempArenaInfo,
	accountInfo data.WGAccountInfo,
	accountList data.WGAccountList,
	clans data.Clans,
	allPlayerShipsStats data.AllPlayerShipsStats,
	allPlayerShipsBadges data.AllPlayerShipsBadges,
) data.Battle {
	friends := make(data.Players, 0)
	enemies := make(data.Players, 0)

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clans[accountID]

		warship, ok := prefetchResult.Warships[vehicle.ShipID]
		if !ok {
			warship = *data.NewUnknownWarship()
		}

		stats := data.NewPersonalStats(
			vehicle.ShipID,
			accountInfo.Data[accountID],
			allPlayerShipsStats.Player(accountID),
			allPlayerShipsBadges[accountID],
			prefetchResult.Warships,
			tempArenaInfo,
		)

		player := data.Player{
			PlayerInfo: data.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo.Data[accountID].HiddenProfile,
			},
			Warship: warship,
			PvPSolo: data.BuildPlayerStats(
				data.StatsPatternPvPSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				prefetchResult.Warships,
			),
			PvPAll: data.BuildPlayerStats(
				data.StatsPatternPvPAll,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				prefetchResult.Warships,
			),
			RankSolo: data.BuildPlayerStats(
				data.StatsPatternRankSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				prefetchResult.Warships,
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
			Arena:    tempArenaInfo.BattleArena(prefetchResult.BattleArenas),
			Type:     tempArenaInfo.BattleType(prefetchResult.BattleTypes),
		},
		Teams: teams,
	}

	return battle
}
