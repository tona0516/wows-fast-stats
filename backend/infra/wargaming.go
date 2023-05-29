package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"fmt"
	"strconv"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

type Wargaming struct {
	config vo.WGConfig
	AppID  string

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
		config:                 config,
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
	w.AppID = appid
}

func (w *Wargaming) AccountInfo(accountIDs []int) (vo.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	return request(
		w.accountInfoClient,
		map[string]string{
			"application_id": w.AppID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         vo.WGAccountInfoData{}.Field(),
			// "extra":          "statistics.pvp_solo",
		},
		apperr.Wg.AccountInfo,
	)
}

func (w *Wargaming) AccountList(accountNames []string) (vo.WGAccountList, error) {
	return request(
		w.accountListClient,
		map[string]string{
			"application_id": w.AppID,
			"search":         strings.Join(accountNames, ","),
			"fields":         vo.WGAccountListData{}.Field(),
			"type":           "exact",
		},
		apperr.Wg.AccountList,
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
			"application_id": w.AppID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         vo.WGClansAccountInfoData{}.Field(),
		},
		apperr.Wg.ClansAccountInfo,
	)
}

func (w *Wargaming) ClansInfo(clanIDs []int) (vo.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	return request(
		w.clansInfoClient,
		map[string]string{
			"application_id": w.AppID,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         vo.WGClansInfoData{}.Field(),
		},
		apperr.Wg.ClansInfo,
	)
}

func (w *Wargaming) ShipsStats(accountID int) (vo.WGShipsStats, error) {
	return request(
		w.shipsStatsClient,
		map[string]string{
			"application_id": w.AppID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         vo.WGShipsStatsData{}.Field(),
			// "extra":          "pvp_solo",
		},
		apperr.Wg.ShipsStats,
	)
}

func (w *Wargaming) EncycShips(pageNo int) (vo.WGEncycShips, error) {
	return request(
		w.encycShipsClient,
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGEncyclopediaShipsData{}.Field(),
			"language":       "en",
			"page_no":        strconv.Itoa(pageNo),
		},
		apperr.Wg.EncyclopediaShips,
	)
}

func (w *Wargaming) EncycInfo() (vo.WGEncycInfo, error) {
	return request(
		w.encycInfoClient,
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGEncyclopediaInfoData{}.Field(),
		},
		apperr.Wg.EncyclopediaInfo,
	)
}

func (w *Wargaming) BattleArenas() (vo.WGBattleArenas, error) {
	return request(
		w.battleArenasClient,
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGBattleArenasData{}.Field(),
			"language":       "ja",
		},
		apperr.Wg.BattleArenas,
	)
}

func (w *Wargaming) BattleTypes() (vo.WGBattleTypes, error) {
	return request(
		w.battleTypesClient,
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGBattleTypesData{}.Field(),
			"language":       "ja",
		},
		apperr.Wg.BattleTypes,
	)
}

func request[T vo.WGResponse](
	client APIClientInterface[T],
	query map[string]string,
	errDetail apperr.AppError,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	operation := func() (T, error) {
		res, err := client.GetRequest(query)

		// Note:
		// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
		message := res.GetError().Message
		if message == "REQUEST_LIMIT_EXCEEDED" || message == "SOURCE_NOT_AVAILABLE" {
			//nolint:goerr113
			return res, fmt.Errorf(message)
		}

		return res, backoff.Permanent(err)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res, errors.WithStack(errDetail.WithRaw(err))
	}
	if res.GetStatus() == "error" {
		//nolint:goerr113
		return res, errors.WithStack(errDetail.WithRaw(fmt.Errorf(res.GetError().Message)))
	}

	return res, err
}
