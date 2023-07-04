package infra

import (
	"strconv"
	"strings"
	"wfs/backend/apperr"
	"wfs/backend/vo"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type Wargaming struct {
	appid                  string
	accountInfoClient      APIClientInterface[vo.WGAccountInfo]
	accountListClient      APIClientInterface[vo.WGAccountList]
	clansAccountInfoClient APIClientInterface[vo.WGClansAccountInfo]
	clansInfoClient        APIClientInterface[vo.WGClansInfo]
	shipsStatsClient       APIClientInterface[vo.WGShipsStats]
	encycShipsClient       APIClientInterface[vo.WGEncycShips]
	encycInfoClient        APIClientInterface[vo.WGEncycInfo]
	battleArenasClient     APIClientInterface[vo.WGBattleArenas]
	battleTypesClient      APIClientInterface[vo.WGBattleTypes]
}

func NewWargaming(config vo.WGConfig) *Wargaming {
	return &Wargaming{
		accountInfoClient:      NewAPIClient[vo.WGAccountInfo](config.BaseURL + "/wows/account/info/"),
		accountListClient:      NewAPIClient[vo.WGAccountList](config.BaseURL + "/wows/account/list/"),
		clansAccountInfoClient: NewAPIClient[vo.WGClansAccountInfo](config.BaseURL + "/wows/clans/accountinfo/"),
		clansInfoClient:        NewAPIClient[vo.WGClansInfo](config.BaseURL + "/wows/clans/info/"),
		shipsStatsClient:       NewAPIClient[vo.WGShipsStats](config.BaseURL + "/wows/ships/stats/"),
		encycShipsClient:       NewAPIClient[vo.WGEncycShips](config.BaseURL + "/wows/encyclopedia/ships/"),
		encycInfoClient:        NewAPIClient[vo.WGEncycInfo](config.BaseURL + "/wows/encyclopedia/info/"),
		battleArenasClient:     NewAPIClient[vo.WGBattleArenas](config.BaseURL + "/wows/encyclopedia/battlearenas/"),
		battleTypesClient:      NewAPIClient[vo.WGBattleTypes](config.BaseURL + "/wows/encyclopedia/battletypes/"),
	}
}

func (w *Wargaming) SetAppID(appid string) {
	w.appid = appid
}

func (w *Wargaming) AccountInfo(accountIDs []int) (vo.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	return request(
		w.accountInfoClient,
		map[string]string{
			"application_id": w.appid,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         vo.WGAccountInfoData{}.Field(),
			"extra":          "statistics.pvp_solo",
		},
	)
}

func (w *Wargaming) AccountList(accountNames []string) (vo.WGAccountList, error) {
	return request(
		w.accountListClient,
		map[string]string{
			"application_id": w.appid,
			"search":         strings.Join(accountNames, ","),
			"fields":         vo.WGAccountListData{}.Field(),
			"type":           "exact",
		},
	)
}

func (w *Wargaming) AccountListForSearch(prefix string) (vo.WGAccountList, error) {
	return request(
		w.accountListClient,
		map[string]string{
			"application_id": w.appid,
			"search":         prefix,
			"fields":         vo.WGAccountListData{}.Field(),
		},
	)
}

func (w *Wargaming) ClansAccountInfo(accountIDs []int) (vo.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	return request(
		w.clansAccountInfoClient,
		map[string]string{
			"application_id": w.appid,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         vo.WGClansAccountInfoData{}.Field(),
		},
	)
}

func (w *Wargaming) ClansInfo(clanIDs []int) (vo.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return vo.WGClansInfo{}, nil
	}

	return request(
		w.clansInfoClient,
		map[string]string{
			"application_id": w.appid,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         vo.WGClansInfoData{}.Field(),
		},
	)
}

func (w *Wargaming) ShipsStats(accountID int) (vo.WGShipsStats, error) {
	return request(
		w.shipsStatsClient,
		map[string]string{
			"application_id": w.appid,
			"account_id":     strconv.Itoa(accountID),
			"fields":         vo.WGShipsStatsData{}.Field(),
			"extra":          "pvp_solo",
		},
	)
}

func (w *Wargaming) EncycShips(pageNo int) (vo.WGEncycShips, error) {
	return request(
		w.encycShipsClient,
		map[string]string{
			"application_id": w.appid,
			"fields":         vo.WGEncyclopediaShipsData{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
	)
}

func (w *Wargaming) EncycInfo() (vo.WGEncycInfo, error) {
	return request(
		w.encycInfoClient,
		map[string]string{
			"application_id": w.appid,
			"fields":         vo.WGEncyclopediaInfoData{}.Field(),
		},
	)
}

func (w *Wargaming) BattleArenas() (vo.WGBattleArenas, error) {
	return request(
		w.battleArenasClient,
		map[string]string{
			"application_id": w.appid,
			"fields":         vo.WGBattleArenasData{}.Field(),
			"language":       "ja",
		},
	)
}

func (w *Wargaming) BattleTypes() (vo.WGBattleTypes, error) {
	return request(
		w.battleTypesClient,
		map[string]string{
			"application_id": w.appid,
			"fields":         vo.WGBattleTypesData{}.Field(),
			"language":       "ja",
		},
	)
}

func request[T vo.WGResponse](
	client APIClientInterface[T],
	query map[string]string,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	operation := func() (APIResponse[T], error) {
		res, err := client.GetRequest(query)

		// Note:
		// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
		message := res.Body.GetError().Message
		if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
			return res, apperr.New(apperr.WargamingAPITemporaryUnavaillalble, errors.New(res.BodyString))
		}

		return res, backoff.Permanent(err)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res.Body, err
	}

	if res.Body.GetStatus() == "error" {
		return res.Body, apperr.New(apperr.WargamingAPIError, errors.New(res.BodyString))
	}

	return res.Body, err
}
