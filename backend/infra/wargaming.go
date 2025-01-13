package infra

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra/response"
	"wfs/backend/infra/webapi"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
	"go.uber.org/ratelimit"
)

type Wargaming struct {
	config RequestConfig
	rl     ratelimit.Limiter
	appID  string
}

func NewWargaming(config RequestConfig, rl ratelimit.Limiter, appID string) *Wargaming {
	return &Wargaming{
		config: config,
		rl:     rl,
		appID:  appID,
	}
}

func (w *Wargaming) AccountInfo(accountIDs []int) (data.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGAccountInfo](
		w,
		"/wows/account/info/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGAccountInfo{}.Field(),
			"extra": strings.Join([]string{
				"statistics.pvp_solo",
				"statistics.pvp_div2",
				"statistics.pvp_div3",
				"statistics.rank_solo",
			}, ","),
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountList(accountNames []string) (data.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": w.appID,
			"search":         strings.Join(accountNames, ","),
			"fields":         response.WGAccountList{}.Field(),
			"type":           "exact",
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountListForSearch(prefix string) (data.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": w.appID,
			"search":         prefix,
			"fields":         response.WGAccountList{}.Field(),
			"limit":          "10",
		},
	)

	return res.Data, err
}

func (w *Wargaming) clansAccountInfo(accountIDs []int) (response.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGClansAccountInfoResponse](
		w,
		"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGClansAccountInfoResponse{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) clansInfo(clanIDs []int) (response.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return response.WGClansInfo{}, nil
	}

	res, err := request[response.WGClansInfoResponse](
		w,
		"/wows/clans/info/",
		map[string]string{
			"application_id": w.appID,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         response.WGClansInfoResponse{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) ShipsStats(accountID int) (data.WGShipsStats, error) {
	res, err := request[response.WGShipsStats](
		w,
		"/wows/ships/stats/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         response.WGShipsStats{}.Field(),
			"extra": strings.Join([]string{
				"pvp_solo",
				"pvp_div2",
				"pvp_div3",
				"rank_solo",
			}, ","),
		},
	)

	return res.Data, err
}

func (w *Wargaming) encycShips(pageNo int) (response.WGEncycShips, int, error) {
	res, err := request[response.WGEncycShips](
		w,
		"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": w.appID,
			"fields":         response.WGEncycShips{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
	)

	return res, res.Meta.PageTotal, err
}

func (w *Wargaming) BattleArenas() (data.WGBattleArenas, error) {
	res, err := request[response.WGBattleArenas](
		w,
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": w.appID,
			"fields":         response.WGBattleArenas{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleTypes() (data.WGBattleTypes, error) {
	res, err := request[response.WGBattleTypes](
		w,
		"/wows/encyclopedia/battletypes/",
		map[string]string{
			"application_id": w.appID,
			"fields":         response.WGBattleTypes{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func request[T response.WGResponse](
	w *Wargaming,
	path string,
	queries map[string]string,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), w.config.Retry)
	operation := func() (T, error) {
		w.rl.Take()
		var result T

		_, body, err := webapi.NewClient(w.config.URL,
			webapi.WithPath(path),
			webapi.WithQuery(queries),
			webapi.WithTimeout(w.config.Timeout),
		).GET()
		if err != nil {
			return result, err
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return result, failure.Translate(err, apperr.WGAPIError)
		}

		if result.GetStatus() == "error" {
			// Note:
			// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
			message := result.GetError().Message
			if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
				return result, failure.New(apperr.WGAPITemporaryUnavaillalble)
			}

			return result, backoff.Permanent(failure.New(apperr.WGAPIError))
		}

		return result, nil
	}

	res, err := backoff.RetryWithData(operation, b)

	return res, err
}
