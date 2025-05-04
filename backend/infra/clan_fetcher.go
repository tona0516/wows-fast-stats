package infra

import (
	"encoding/json"
	"regexp"
	"slices"
	"strings"
	"sync"
	"wfs/backend/domain/model"
	"wfs/backend/infra/webapi"

	"github.com/abadojack/whatlanggo"
	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

var urlRegExp = regexp.MustCompile(`https?://[^\s]+`)

type ClanFetcher struct {
	wargaming                 webapi.Wargaming
	unofficialWargamingClient req.Client
}

func NewClanFetcher(
	wargaming webapi.Wargaming,
	unofficialWargamingClient req.Client,
) *ClanFetcher {
	return &ClanFetcher{
		wargaming:                 wargaming,
		unofficialWargamingClient: unofficialWargamingClient,
	}
}

func (f *ClanFetcher) Fetch(accountIDs []int) (model.Clans, error) {
	clansAccountInfo, err := f.clansAccountInfo(accountIDs)
	if err != nil {
		return nil, err
	}

	clans, err := f.clanInfo(clansAccountInfo)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0, len(clans))
	for _, clan := range clans {
		if slices.Contains(tags, clan.Tag) {
			continue
		}
		tags = append(tags, clan.Tag)
	}

	hexColors, err := f.hexColor(tags)
	if err != nil {
		return nil, err
	}

	for accountID, clan := range clans {
		clans[accountID] = model.Clan{
			ID:          clan.ID,
			Tag:         clan.Tag,
			Description: clan.Description,
			HexColor:    hexColors[clan.Tag],
			Lang:        f.lang(clan.Description),
		}
	}

	return clans, nil
}

func (f *ClanFetcher) clansAccountInfo(accountIDs []int) (map[int]int, error) {
	result := make(map[int]int, 0)

	clansAccountInfo, err := f.wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		return nil, err
	}

	for accountID, clan := range clansAccountInfo {
		if clan.ClanID == 0 {
			continue
		}
		result[accountID] = clan.ClanID
	}

	return result, nil
}

func (f *ClanFetcher) clanInfo(clanIDMap map[int]int) (model.Clans, error) {
	result := make(model.Clans)

	clanIDs := make([]int, 0, len(clanIDMap))
	for _, clanID := range clanIDMap {
		clanIDs = append(clanIDs, clanID)
	}

	clansInfo, err := f.wargaming.ClansInfo(clanIDs)
	if err != nil {
		return nil, err
	}

	for accountID, clanID := range clanIDMap {
		clanInfo, ok := clansInfo[clanID]
		if !ok {
			continue
		}

		result[accountID] = model.Clan{
			ID:          clanID,
			Tag:         clanInfo.Tag,
			Description: clanInfo.Description,
		}
	}

	return result, nil
}

func (f *ClanFetcher) hexColor(tags []string) (map[string]string, error) {
	result := make(map[string]string)

	var mu sync.Mutex
	err := doParallel(tags, func(tag string) error {
		resp, err := f.unofficialWargamingClient.R().SetQueryParams(map[string]string{
			"search": tag,
			"type":   "clans",
		}).Get("/api/search/autocomplete/")
		if err != nil {
			return failure.Wrap(err)
		}

		var autocomplete UWGClansAutocomplete
		if err := json.Unmarshal(resp.Bytes(), &autocomplete); err != nil {
			return failure.Wrap(err)
		}

		hexColor := autocomplete.HexColor(tag)
		if hexColor != "" {
			mu.Lock()
			result[tag] = hexColor
			mu.Unlock()
		}

		return nil
	})

	return result, err
}

func (f *ClanFetcher) lang(description string) string {
	options := whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Jpn: true,
			whatlanggo.Kor: true,
			whatlanggo.Cmn: true,
		},
	}

	description = urlRegExp.ReplaceAllString(description, "")
	description = strings.ReplaceAll(description, "\n", "")

	if len(description) == 0 {
		return ""
	}

	info := whatlanggo.DetectWithOptions(description, options)

	return info.Lang.Iso6391()
}
