package infra

import (
	"encoding/json"
	"strconv"
	"strings"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra/response"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
	"go.uber.org/ratelimit"
	"golang.org/x/exp/slices"
)

type Wargaming struct {
	config  apiConfig
	limiter ratelimit.Limiter
}

func NewWargaming(config apiConfig, limiter ratelimit.Limiter) *Wargaming {
	return &Wargaming{config: config, limiter: limiter}
}

func (w *Wargaming) AccountInfo(appID string, accountIDs []int) (data.WGAccountInfo, error) {
	if len(accountIDs) == 0 {
		return data.WGAccountInfo{}, nil
	}

	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGAccountInfo](
		"/wows/account/info/",
		map[string]string{
			"application_id": appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGAccountInfo{}.Field(),
			"extra":          "statistics.pvp_solo,statistics.pvp_div2,statistics.pvp_div3",
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) AccountList(appID string, accountNames []string) (data.WGAccountList, error) {
	if len(accountNames) == 0 {
		return data.WGAccountList{}, nil
	}

	res, err := request[response.WGAccountList](
		"/wows/account/list/",
		map[string]string{
			"application_id": appID,
			"search":         strings.Join(accountNames, ","),
			"fields":         response.WGAccountList{}.Field(),
			"type":           "exact",
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) AccountListForSearch(appID string, prefix string) (data.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		"/wows/account/list/",
		map[string]string{
			"application_id": appID,
			"search":         prefix,
			"fields":         response.WGAccountList{}.Field(),
			"limit":          "10",
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) ClansAccountInfo(appID string, accountIDs []int) (data.WGClansAccountInfo, error) {
	if len(accountIDs) == 0 {
		return data.WGClansAccountInfo{}, nil
	}

	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGClansAccountInfo](
		"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGClansAccountInfo{}.Field(),
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) ClansInfo(appID string, clanIDs []int) (data.WGClansInfo, error) {
	if len(clanIDs) == 0 {
		return data.WGClansInfo{}, nil
	}

	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGClansInfo](
		"/wows/clans/info/",
		map[string]string{
			"application_id": appID,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         response.WGClansInfo{}.Field(),
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) ShipsStats(appID string, accountID int) (data.WGShipsStats, error) {
	res, err := request[response.WGShipsStats](
		"/wows/ships/stats/",
		map[string]string{
			"application_id": appID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         response.WGShipsStats{}.Field(),
			"extra":          "pvp_solo,pvp_div2,pvp_div3",
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) EncycShips(appID string, pageNo int) (data.WGEncycShips, int, error) {
	res, err := request[response.WGEncycShips](
		"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycShips{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
		w.config,
		w.limiter,
	)

	return res.Data, res.Meta.PageTotal, err
}

func (w *Wargaming) EncycInfo(appID string) (data.WGEncycInfoData, error) {
	res, err := request[response.WGEncycInfo](
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycInfo{}.Field(),
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) BattleArenas(appID string) (data.WGBattleArenas, error) {
	res, err := request[response.WGBattleArenas](
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGBattleArenas{}.Field(),
			"language":       "ja",
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) BattleTypes(appID string) (data.WGBattleTypes, error) {
	res, err := request[response.WGBattleTypes](
		"/wows/encyclopedia/battletypes/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGBattleTypes{}.Field(),
			"language":       "ja",
		},
		w.config,
		w.limiter,
	)

	return res.Data, err
}

func (w *Wargaming) Test(appID string) bool {
	_, err := request[response.WGEncycInfo](
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": appID,
			"fields":         response.WGEncycInfo{}.Field(),
		},
		w.config,
		w.limiter,
	)

	return err == nil
}

func getError(r *resty.Response, err error) (error, bool) {
	errCtx := failure.Context{
		"url":         r.Request.URL,
		"status_code": strconv.Itoa(r.StatusCode()),
		"body":        string(r.Body()),
	}

	if err != nil {
		return failure.New(apperr.WGAPIError, failure.Messagef(err.Error()), errCtx), true
	}

	body := response.WGResponse[any]{}
	if err := json.Unmarshal(r.Body(), &body); err != nil {
		return failure.New(apperr.WGAPIParseResponseError, failure.Messagef(err.Error()), errCtx), false
	}

	if body.Status != "error" {
		return nil, false
	}

	msg := body.Error.Message
	// Note:
	// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
	if slices.Contains(
		[]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"},
		msg,
	) {
		return failure.New(apperr.WGAPITemporaryUnavaillalble, failure.Messagef(msg), errCtx), true
	}

	return failure.New(apperr.WGAPIError, failure.Messagef(msg), errCtx), false
}

func request[T any](
	path string,
	queries map[string]string,
	config apiConfig,
	limiter ratelimit.Limiter,
) (T, error) {
	limiter.Take()

	client := resty.New().
		SetTimeout(config.timeout).
		SetRetryCount(config.retryCount)

	var result T
	r, err := client.R().
		SetResult(&result).
		SetQueryParams(queries).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			_, shouldRetry := getError(r, err)
			return shouldRetry
		}).
		Get(config.baseURL + path)

	err, _ = getError(r, err)
	return result, err
}
