package webapi

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"wfs/backend/apperr"
	"wfs/backend/infra/response"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
	"go.uber.org/ratelimit"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type Wargaming interface {
	AccountInfo(accountIDs []int) (response.WGAccountInfo, error)
	AccountList(accountNames []string) (response.WGAccountList, error)
	AccountListForSearch(prefix string) (response.WGAccountList, error)
	ClansAccountInfo(accountIDs []int) (response.WGClansAccountInfo, error)
	ClansInfo(clanIDs []int) (response.WGClansInfo, error)
	ShipsStats(accountID int) (response.WGShipsStats, error)
	EncycShips(pageNo int) (response.WGEncycShips, int, error)
	BattleArenas() (response.WGBattleArenas, error)
	BattleTypes() (response.WGBattleTypes, error)
	GameVersion() (string, error)
}

type wargaming struct {
	config RequestConfig
	rl     ratelimit.Limiter
	appID  string
}

func NewWargaming(config RequestConfig, rl ratelimit.Limiter, appID string) Wargaming {
	return &wargaming{
		config: config,
		rl:     rl,
		appID:  appID,
	}
}

func (w *wargaming) AccountInfo(accountIDs []int) (response.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGAccountInfoResponse](
		w,
		"/wows/account/info/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGAccountInfoResponse{}.Field(),
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

func (w *wargaming) AccountList(accountNames []string) (response.WGAccountList, error) {
	res, err := request[response.WGAccountListResponse](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": w.appID,
			"search":         strings.Join(accountNames, ","),
			"fields":         response.WGAccountListResponse{}.Field(),
			"type":           "exact",
		},
	)

	return res.Data, err
}

func (w *wargaming) AccountListForSearch(prefix string) (response.WGAccountList, error) {
	res, err := request[response.WGAccountListResponse](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": w.appID,
			"search":         prefix,
			"fields":         response.WGAccountListResponse{}.Field(),
			"limit":          "10",
		},
	)

	return res.Data, err
}

func (w *wargaming) ClansAccountInfo(accountIDs []int) (response.WGClansAccountInfo, error) {
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

func (w *wargaming) ClansInfo(clanIDs []int) (response.WGClansInfo, error) {
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

func (w *wargaming) ShipsStats(accountID int) (response.WGShipsStats, error) {
	res, err := request[response.WGShipsStatsResponse](
		w,
		"/wows/ships/stats/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         response.WGShipsStatsResponse{}.Field(),
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

func (w *wargaming) EncycShips(pageNo int) (response.WGEncycShips, int, error) {
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

func (w *wargaming) BattleArenas() (response.WGBattleArenas, error) {
	res, err := request[response.WGBattleArenasResponse](
		w,
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": w.appID,
			"fields":         response.WGBattleArenasResponse{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *wargaming) BattleTypes() (response.WGBattleTypes, error) {
	res, err := request[response.WGBattleTypesResponse](
		w,
		"/wows/encyclopedia/battletypes/",
		map[string]string{
			"application_id": w.appID,
			"fields":         response.WGBattleTypesResponse{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *wargaming) GameVersion() (string, error) {
	res, err := request[response.WGEncycInfo](
		w,
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": w.appID,
			"fields":         "game_version",
		},
	)

	return res.Data.GameVersion, err
}

func request[T response.WGResponse](
	w *wargaming,
	path string,
	queries map[string]string,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), w.config.Retry)
	operation := func() (T, error) {
		w.rl.Take()
		var result T

		_, body, err := NewClient(w.config.URL,
			WithPath(path),
			WithQuery(queries),
			WithTimeout(w.config.Timeout),
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
