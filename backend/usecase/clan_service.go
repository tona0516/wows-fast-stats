package usecase

import (
	"context"
	"strings"
	"sync"
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/abadojack/whatlanggo"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

type clanService struct {
	wargamingClient adapter.WargamingClient
	clanClient      adapter.ClanClient
}

func NewClanService(i do.Injector) (*clanService, error) {
	return &clanService{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		clanClient:      do.MustInvoke[adapter.ClanClient](i),
	}, nil
}

func (s *clanService) fetchAll(
	ctx context.Context,
	accountIDs []data.AccountID,
) (data.Clans, error) {
	result := make(data.Clans)

	clansAccountInfo, err := s.wargamingClient.ClansAccountInfo(ctx, accountIDs)
	if err != nil {
		return nil, err
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := s.wargamingClient.ClansInfo(ctx, clanIDs)
	if err != nil {
		return nil, err
	}

	colorMap, err := s.fetchClanColor(ctx, clansInfo)
	if err != nil {
		return nil, err
	}

	languageMap := s.clanLanguage(clansInfo)

	for _, accountID := range accountIDs {
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		hexColor := colorMap[clanID]
		language := languageMap[clanID]

		result[accountID] = data.Clan{
			ID:       clanID,
			Tag:      clanTag,
			HexColor: hexColor,
			Language: language,
		}
	}

	return result, nil
}

func (s *clanService) fetchClanColor(
	ctx context.Context,
	clansInfo data.WGClansInfo,
) (map[data.ClanID]string, error) {
	result := make(map[data.ClanID]string)
	eg, egCtx := errgroup.WithContext(ctx)

	var mu sync.Mutex
	for clanID, clanInfo := range clansInfo.Data {
		eg.Go(func() error {
			autocomplete, err := s.clanClient.ClanAutoComplete(egCtx, clanInfo.Tag)
			if err != nil {
				return err
			}

			hexColor := autocomplete.HexColor(clanID)
			if hexColor == "" {
				return nil
			}

			mu.Lock()
			result[clanID] = hexColor
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *clanService) clanLanguage(clansInfo data.WGClansInfo) map[data.ClanID]string {
	result := make(map[data.ClanID]string)

	options := whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Jpn: true,
			whatlanggo.Kor: true,
			whatlanggo.Cmn: true,
		},
	}

	for clanID, clanInfo := range clansInfo.Data {
		// URLを空文字に
		description := urlRegex.ReplaceAllString(clanInfo.Description, "")
		// 改行を空文字に
		description = strings.ReplaceAll(description, "\n", "")

		if len(description) == 0 {
			continue
		}

		info := whatlanggo.DetectWithOptions(description, options)
		result[clanID] = info.Lang.Iso6391()
	}

	return result
}
