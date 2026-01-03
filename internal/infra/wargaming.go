package infra

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"time"
	"wfs/internal/apperr"
	"wfs/internal/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"go.uber.org/ratelimit"
)

type Wargaming struct {
	url           string
	maxRetry      int
	timeoutSec    int
	retryInterval int
	limier        ratelimit.Limiter
	appID         string
}

func NewWargaming(
	url string,
	maxRetry int,
	timeoutSec int,
	retryInterval int,
	rateLimitRPS int,
	appID string,
) *Wargaming {
	return &Wargaming{
		url:           url,
		maxRetry:      maxRetry,
		timeoutSec:    timeoutSec,
		retryInterval: retryInterval,
		limier:        ratelimit.New(rateLimitRPS),
		appID:         appID,
	}
}

func (w *Wargaming) AccountInfo(accountIDs []int) (data.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[data.WGAccountInfo](
		w,
		"/wows/account/info/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         data.WGAccountInfo{}.Field(),
			"extra": strings.Join([]string{
				"statistics.pvp_solo",
				"statistics.pvp_div2",
				"statistics.pvp_div3",
				"statistics.rank_solo",
			}, ","),
		},
	)

	return res, err
}

func (w *Wargaming) AccountList(accountNames []string) (data.WGAccountList, error) {
	res, err := request[data.WGAccountList](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": w.appID,
			"search":         strings.Join(accountNames, ","),
			"fields":         data.WGAccountList{}.Field(),
			"type":           "exact",
		},
	)

	return res, err
}

func (w *Wargaming) AccountListForSearch(prefix string) (data.WGAccountList, error) {
	res, err := request[data.WGAccountList](
		w,
		"/wows/account/list/",
		map[string]string{
			"application_id": w.appID,
			"search":         prefix,
			"fields":         data.WGAccountList{}.Field(),
			"limit":          "10",
		},
	)

	return res, err
}

func (w *Wargaming) ClansAccountInfo(accountIDs []int) (data.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[data.WGClansAccountInfo](
		w,
		"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         data.WGClansAccountInfo{}.Field(),
		},
	)

	return res, err
}

func (w *Wargaming) ClansInfo(clanIDs []int) (data.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return data.WGClansInfo{}, nil
	}

	res, err := request[data.WGClansInfo](
		w,
		"/wows/clans/info/",
		map[string]string{
			"application_id": w.appID,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         data.WGClansInfo{}.Field(),
		},
	)

	return res, err
}

func (w *Wargaming) ShipsStats(accountID int) (data.WGShipsStats, error) {
	res, err := request[data.WGShipsStats](
		w,
		"/wows/ships/stats/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         data.WGShipsStats{}.Field(),
			"extra": strings.Join([]string{
				"pvp_solo",
				"pvp_div2",
				"pvp_div3",
				"rank_solo",
			}, ","),
		},
	)

	return res, err
}

func (w *Wargaming) EncycShips(pageNo int) (data.WGEncycShips, error) {
	res, err := request[data.WGEncycShips](
		w,
		"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": w.appID,
			"fields":         data.WGEncycShips{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
	)

	return res, err
}

func (w *Wargaming) EncycInfo() (data.WGEncycInfoData, error) {
	res, err := request[data.WGEncycInfo](
		w,
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": w.appID,
			"fields":         data.WGEncycInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleArenas() (data.WGBattleArenas, error) {
	res, err := request[data.WGBattleArenas](
		w,
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": w.appID,
			"fields":         data.WGBattleArenas{}.Field(),
			"language":       "ja",
		},
	)

	return res, err
}

func (w *Wargaming) BattleTypes() (data.WGBattleTypes, error) {
	res, err := request[data.WGBattleTypes](
		w,
		"/wows/encyclopedia/battletypes/",
		map[string]string{
			"application_id": w.appID,
			"fields":         data.WGBattleTypes{}.Field(),
			"language":       "ja",
		},
	)

	return res, err
}

func (w *Wargaming) ShipsBadges(accountID int) (data.WGShipsBadges, error) {
	res, err := request[data.WGShipsBadges](
		w,
		"/wows/ships/badges/",
		map[string]string{
			"application_id": w.appID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         data.WGShipsBadges{}.Field(),
		},
	)

	return res, err
}

func request[T data.WGResponse](
	w *Wargaming,
	path string,
	queries map[string]string,
) (T, error) {
	c := req.C().
		SetBaseURL(w.url).
		SetTimeout(time.Duration(w.timeoutSec) * time.Second).
		SetCommonRetryCount(w.maxRetry).
		SetCommonRetryFixedInterval(time.Duration(w.retryInterval) * time.Millisecond).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}

			var body data.WGResponseCommon[any]
			if err := json.Unmarshal(resp.Bytes(), &body); err == nil {
				err := convertError(body.Status, body.Error.Message)
				if failure.Is(err, apperr.WGAPITemporaryUnavaillalble) {
					return true
				}
			}

			return false
		}).
		OnBeforeRequest(func(client *req.Client, req *req.Request) error {
			w.limier.Take()
			return nil
		}).
		SetCommonRetryHook(func(resp *req.Response, err error) {
			w.limier.Take()
		})

	var result T
	_, err := c.R().
		SetSuccessResult(&result).
		SetQueryParams(queries).
		Get(path)

	if err != nil {
		return result, failure.Translate(err, apperr.WGAPIError)
	}

	if err := convertError(result.GetStatus(), result.GetError().Message); err != nil {
		return result, failure.Wrap(err)
	}

	return result, nil
}

func convertError(status string, message string) error {
	if status == "error" {
		// Note:
		// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
		if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
			return failure.New(apperr.WGAPITemporaryUnavaillalble)
		}

		return failure.New(apperr.WGAPIError)
	}

	return nil
}
