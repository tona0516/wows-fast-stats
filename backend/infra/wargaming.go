package infra

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra/response"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
	"go.uber.org/ratelimit"
	"golang.org/x/exp/slices"
)

type Wargaming struct {
	baseURL string
	rl      ratelimit.Limiter
}

func NewWargaming(baseURL string, rl ratelimit.Limiter) *Wargaming {
	return &Wargaming{
		baseURL: baseURL,
		rl:      rl,
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

func convertWGError(body []byte) error {
	resp := response.WGResponseCommon[any]{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return failure.Wrap(err)
	}

	if resp.Status == "error" {
		// Note:
		// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
		message := resp.Error.Message
		if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
			return failure.New(apperr.WGAPITemporaryUnavaillalble)
		}

		return failure.New(apperr.WGAPIError)
	}

	return nil
}

func request[T response.WGResponse](
	w *Wargaming,
	path string,
	queries map[string]string,
) (T, error) {
	w.rl.Take()

	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(2).
		AddRetryCondition(func(resp *resty.Response, err error) bool {
			if err != nil {
				return true
			}

			if err := convertWGError(resp.Body()); failure.Is(err, apperr.WGAPITemporaryUnavaillalble) {
				return true
			}

			return false
		})

	var result T
	resp, err := client.R().
		SetResult(&result).
		SetQueryParams(queries).
		Get(w.baseURL + path)

	errCtx := failure.Context{
		"url":         resp.Request.URL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}
	if err != nil {
		return result, failure.Wrap(err, errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return result, failure.New(apperr.GithubAPICheckUpdateError, errCtx)
	}

	if err := convertWGError(resp.Body()); err != nil {
		return result, failure.Wrap(err, errCtx)
	}

	return result, nil
}
