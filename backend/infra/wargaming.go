package infra

import (
	"strconv"
	"strings"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/infra/response"
	"wfs/backend/infra/webapi"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slices"
)

type Wargaming struct {
	config RequestConfig
}

func NewWargaming(config RequestConfig) *Wargaming {
	return &Wargaming{config: config}
}

func (w *Wargaming) AccountInfo(appID string, accountIDs []int) (model.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGAccountInfo](
		w.config.URL+"/wows/account/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGAccountInfo{}.Field(),
			"extra":          "statistics.pvp_solo,statistics.pvp_div2,statistics.pvp_div3",
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountList(appID string, accountNames []string) (model.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w.config.URL+"/wows/account/list/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"search":         strings.Join(accountNames, ","),
			"fields":         response.WGAccountList{}.Field(),
			"type":           "exact",
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountListForSearch(appID string, prefix string) (model.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w.config.URL+"/wows/account/list/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"search":         prefix,
			"fields":         response.WGAccountList{}.Field(),
			"limit":          "10",
		},
	)

	return res.Data, err
}

func (w *Wargaming) ClansAccountInfo(appID string, accountIDs []int) (model.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGClansAccountInfo](
		w.config.URL+"/wows/clans/accountinfo/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGClansAccountInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) ClansInfo(appID string, clanIDs []int) (model.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return model.WGClansInfo{}, nil
	}

	res, err := request[response.WGClansInfo](
		w.config.URL+"/wows/clans/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         response.WGClansInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) ShipsStats(appID string, accountID int) (model.WGShipsStats, error) {
	res, err := request[response.WGShipsStats](
		w.config.URL+"/wows/ships/stats/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         response.WGShipsStats{}.Field(),
			"extra":          "pvp_solo,pvp_div2,pvp_div3",
		},
	)

	return res.Data, err
}

func (w *Wargaming) EncycShips(appID string, pageNo int) (model.WGEncycShips, int, error) {
	res, err := request[response.WGEncycShips](
		w.config.URL+"/wows/encyclopedia/ships/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycShips{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
	)

	return res.Data, res.Meta.PageTotal, err
}

func (w *Wargaming) EncycInfo(appID string) (model.WGEncycInfoData, error) {
	res, err := request[response.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleArenas(appID string) (model.WGBattleArenas, error) {
	res, err := request[response.WGBattleArenas](
		w.config.URL+"/wows/encyclopedia/battlearenas/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"fields":         response.WGBattleArenas{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleTypes(appID string) (model.WGBattleTypes, error) {
	res, err := request[response.WGBattleTypes](
		w.config.URL+"/wows/encyclopedia/battletypes/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"fields":         response.WGBattleTypes{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *Wargaming) Test(appID string) (bool, error) {
	_, err := request[response.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycInfo{}.Field(),
		},
	)

	return err == nil, err
}

func request[T response.WGResponse](
	rawURL string,
	retry uint64,
	timeout time.Duration,
	queries map[string]string,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retry)
	operation := func() (webapi.Response[any, T], error) {
		res, err := webapi.GetRequest[T](rawURL, timeout, queries)
		errCtx := failure.Context{
			"url":         res.Request.URL,
			"status_code": strconv.Itoa(res.StatusCode),
			"body":        string(res.BodyByte),
		}

		if err != nil {
			return res, failure.Wrap(err, errCtx)
		}

		if res.Body.GetStatus() == "error" {
			// Note:
			// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
			message := res.Body.GetError().Message
			if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
				return res, failure.New(apperr.WGAPITemporaryUnavaillalble, errCtx)
			}

			return res, backoff.Permanent(failure.New(apperr.WGAPIError, errCtx))
		}

		return res, nil
	}
	res, err := backoff.RetryWithData(operation, b)

	return res.Body, failure.Wrap(err)
}
