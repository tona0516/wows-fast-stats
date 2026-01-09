package service

import (
	"regexp"
	"strings"
	"sync"
	"wfs/internal/adapter"
	"wfs/internal/data"

	"github.com/abadojack/whatlanggo"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

// URLを検出する正規表現パターン.
var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

type UserData struct {
	AccountInfo          data.WGAccountInfo
	AccountList          data.WGAccountList
	Clans                data.Clans
	AllPlayerShipsStats  data.AllPlayerShipsStats
	AllPlayerShipsBadges data.AllPlayerShipsBadges
}

type UserDataFetcher struct {
	wargamingClient adapter.WargamingClient
	clanClient      adapter.ClanClient
}

func NewUserDataFetcher(i do.Injector) (*UserDataFetcher, error) {
	return &UserDataFetcher{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		clanClient:      do.MustInvoke[adapter.ClanClient](i),
	}, nil
}

func (f *UserDataFetcher) Fetch(accountNames []string) (*UserData, error) {
	accountList, err := f.wargamingClient.AccountList(accountNames)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	accountIDs := accountList.AccountIDs()

	eg := errgroup.Group{}

	var accountInfo data.WGAccountInfo
	eg.Go(func() error {
		var err error
		accountInfo, err = f.wargamingClient.AccountInfo(accountIDs)
		return err
	})

	var allShipStats data.AllPlayerShipsStats
	eg.Go(func() error {
		var err error
		allShipStats, err = f.fetchAllPlayerShipsStats(accountIDs)
		return err
	})

	var allShipBadges data.AllPlayerShipsBadges
	eg.Go(func() error {
		var err error
		allShipBadges, err = f.fetchAllPlayerShipsBadges(accountIDs)
		return err
	})

	var clans data.Clans
	eg.Go(func() error {
		var err error
		clans, err = f.fetchClan(accountIDs)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, failure.Wrap(err)
	}

	return &UserData{
		AccountInfo:          accountInfo,
		AccountList:          accountList,
		Clans:                clans,
		AllPlayerShipsStats:  allShipStats,
		AllPlayerShipsBadges: allShipBadges,
	}, nil
}

func (f *UserDataFetcher) fetchAllPlayerShipsStats(accountIDs []int) (data.AllPlayerShipsStats, error) {
	result := make(data.AllPlayerShipsStats)
	var mu sync.Mutex

	eg := errgroup.Group{}

	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := f.wargamingClient.ShipsStats(accountID)
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

func (f *UserDataFetcher) fetchClan(accountIDs []int) (data.Clans, error) {
	result := make(data.Clans)

	clansAccountInfo, err := f.wargamingClient.ClansAccountInfo(accountIDs)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := f.wargamingClient.ClansInfo(clanIDs)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	clanInfoArray := clansInfo.ToArray()
	colorMap, err := f.fetchClanColor(clanInfoArray)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	languageMap := f.fetchClanLanguage(clanInfoArray)

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

func (f *UserDataFetcher) fetchClanColor(clanInfoSlice []data.WGClansInfoData) (map[string]string, error) {
	result := make(map[string]string)
	var mu sync.Mutex

	eg := errgroup.Group{}
	for _, clanInfo := range clanInfoSlice {
		eg.Go(func() error {
			autocomplete, err := f.clanClient.ClanAutoComplete(clanInfo.Tag)
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

func (f *UserDataFetcher) fetchClanLanguage(clanInfoSilce []data.WGClansInfoData) map[string]string {
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

func (f *UserDataFetcher) fetchAllPlayerShipsBadges(accountIDs []int) (data.AllPlayerShipsBadges, error) {
	result := make(data.AllPlayerShipsBadges)
	var mu sync.Mutex

	eg := errgroup.Group{}
	for _, accountID := range accountIDs {
		eg.Go(func() error {
			shipsBadges, err := f.wargamingClient.ShipsBadges(accountID)
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
