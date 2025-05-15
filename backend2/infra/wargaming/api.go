package wargaming

import (
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/do"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type API interface {
	AccountInfo(accountIDs []int) (AccountInfo, error)
	AccountList(accountNames []string) (AccountList, error)
	AccountListForSearch(prefix string) (AccountList, error)
	ClansAccountInfo(accountIDs []int) (ClansAccountInfo, error)
	ClansInfo(clanIDs []int) (ClansInfo, error)
	ShipsStats(accountID int) (ShipsStats, error)
	EncycShips(pageNo int) (EncycShips, error)
	BattleArenas() (BattleArenas, error)
	BattleTypes() (BattleTypes, error)
	GameVersion() (string, error)
}

type api struct {
	client *req.Client
	appID  string
}

func NewAPI(i *do.Injector) (API, error) {
	return &api{
		client: do.MustInvokeNamed[*req.Client](i, "WargamingAPIClient"),
		appID:  do.MustInvokeNamed[string](i, "WargamingAppID"),
	}, nil
}

func (w *api) AccountInfo(accountIDs []int) (AccountInfo, error) {
	strAccountIDs := toStrings(accountIDs)

	result := AccountInfoResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"account_id":     strings.Join(strAccountIDs, ","),
		"fields":         result.Field(),
		"extra": strings.Join([]string{
			"statistics.pvp_solo",
			"statistics.pvp_div2",
			"statistics.pvp_div3",
			"statistics.rank_solo",
		}, ","),
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/account/info/")

	return result.Data, err
}

func (w *api) AccountList(accountNames []string) (AccountList, error) {
	result := AccountListResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"search":         strings.Join(accountNames, ","),
		"fields":         result.Field(),
		"type":           "exact",
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/account/list/")

	return result.Data, err
}

func (w *api) AccountListForSearch(prefix string) (AccountList, error) {
	result := AccountListResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"search":         prefix,
		"fields":         result.Field(),
		"limit":          "10",
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/account/list/")

	return result.Data, err
}

func (w *api) ClansAccountInfo(accountIDs []int) (ClansAccountInfo, error) {
	strAccountIDs := toStrings(accountIDs)

	result := ClansAccountInfoResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"account_id":     strings.Join(strAccountIDs, ","),
		"fields":         result.Field(),
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/clans/accountinfo/")

	return result.Data, err
}

func (w *api) ClansInfo(clanIDs []int) (ClansInfo, error) {
	strClanIDs := toStrings(clanIDs)

	result := ClansInfoResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"clan_id":        strings.Join(strClanIDs, ","),
		"fields":         result.Field(),
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/clans/info/")

	return result.Data, err
}

func (w *api) ShipsStats(accountID int) (ShipsStats, error) {
	result := ShipsStatsResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"account_id":     strconv.Itoa(accountID),
		"fields":         result.Field(),
		"extra": strings.Join([]string{
			"pvp_solo",
			"pvp_div2",
			"pvp_div3",
			"rank_solo",
		}, ","),
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/ships/stats/")

	return result.Data, err
}

func (w *api) EncycShips(pageNo int) (EncycShips, error) {
	result := EncycShips{}

	params := map[string]string{
		"application_id": w.appID,
		"fields":         result.Field(),
		"language":       "ja",
		"page_no":        strconv.Itoa(pageNo),
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/encyclopedia/ships/")

	return result, err
}

func (w *api) BattleArenas() (BattleArenas, error) {
	result := BattleArenasResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"fields":         result.Field(),
		"language":       "ja",
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/encyclopedia/battlearenas/")

	return result.Data, err
}

func (w *api) BattleTypes() (BattleTypes, error) {
	result := BattleTypesResponse{}

	params := map[string]string{
		"application_id": w.appID,
		"fields":         result.Field(),
		"language":       "ja",
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/encyclopedia/battletypes/")

	return result.Data, err
}

func (w *api) GameVersion() (string, error) {
	result := EncycInfo{}

	params := map[string]string{
		"application_id": w.appID,
		"fields":         "game_version",
	}

	_, err := w.client.R().
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/wows/encyclopedia/info/")

	return result.Data.GameVersion, err
}
