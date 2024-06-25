package infra

import (
	"strconv"
	"strings"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra/response"
	"wfs/backend/infra/webapi"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
	"go.uber.org/ratelimit"
	"golang.org/x/exp/slices"
)

type Wargaming struct {
	config RequestConfig
	rl     ratelimit.Limiter
}

func NewWargaming(config RequestConfig, rl ratelimit.Limiter) *Wargaming {
	return &Wargaming{
		config: config,
		rl:     rl,
	}
}

func (w *Wargaming) AccountInfo(appID string, accountIDs []int) (data.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGAccountInfo](
		w,
		"/wows/account/info/",
		map[string]string{
			"application_id": appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGAccountInfo{}.Field(),
			"extra":          "statistics.pvp_solo,statistics.pvp_div2,statistics.pvp_div3",
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountList(appID string, accountNames []string) (data.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": appID,
			"search":         strings.Join(accountNames, ","),
			"fields":         response.WGAccountList{}.Field(),
			"type":           "exact",
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountListForSearch(appID string, prefix string) (data.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": appID,
			"search":         prefix,
			"fields":         response.WGAccountList{}.Field(),
			"limit":          "10",
		},
	)

	return res.Data, err
}

func (w *Wargaming) ClansAccountInfo(appID string, accountIDs []int) (data.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGClansAccountInfo](
		w,
		"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGClansAccountInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) ClansInfo(appID string, clanIDs []int) (data.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return data.WGClansInfo{}, nil
	}

	res, err := request[response.WGClansInfo](
		w,
		"/wows/clans/info/",
		map[string]string{
			"application_id": appID,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         response.WGClansInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) ShipsStats(appID string, accountID int) (data.WGShipsStats, error) {
	res, err := request[response.WGShipsStats](
		w,
		"/wows/ships/stats/",
		map[string]string{
			"application_id": appID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         response.WGShipsStats{}.Field(),
			"extra":          "pvp_solo,pvp_div2,pvp_div3",
		},
	)

	return res.Data, err
}

func (w *Wargaming) EncycShips(appID string, pageNo int) (data.WGEncycShips, int, error) {
	res, err := request[response.WGEncycShips](
		w,
		"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycShips{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
	)

	return res.Data, res.Meta.PageTotal, err
}

func (w *Wargaming) EncycInfo(appID string) (data.WGEncycInfoData, error) {
	res, err := request[response.WGEncycInfo](
		w,
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleArenas(appID string) (data.WGBattleArenas, error) {
	res, err := request[response.WGBattleArenas](
		w,
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGBattleArenas{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleTypes(appID string) (data.WGBattleTypes, error) {
	res, err := request[response.WGBattleTypes](
		w,
		"/wows/encyclopedia/battletypes/",
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
		w,
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycInfo{}.Field(),
		},
	)

	return err == nil, err
}

func request[T response.WGResponse](
	w *Wargaming,
	path string,
	queries map[string]string,
) (T, error) {
	url := w.config.URL + path
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), w.config.Retry)
	operation := func() (webapi.Response[any, T], error) {
		w.rl.Take()
		res, err := webapi.GetRequest[T](url, w.config.Timeout, queries, w.config.Transport)
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
