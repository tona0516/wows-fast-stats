package infra

import (
	"strconv"
	"strings"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/domain"
	"wfs/backend/infra/response"
	"wfs/backend/infra/webapi"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slices"
)

type Wargaming struct {
	appid  string
	config RequestConfig
}

func NewWargaming(config RequestConfig) *Wargaming {
	return &Wargaming{config: config}
}

func (w *Wargaming) SetAppID(appid string) {
	w.appid = appid
}

func (w *Wargaming) AccountInfo(accountIDs []int) (domain.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGAccountInfo](
		w.config.URL+"/wows/account/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGAccountInfo{}.Field(),
			"extra":          "statistics.pvp_solo,statistics.pvp_div2,statistics.pvp_div3",
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountList(accountNames []string) (domain.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w.config.URL+"/wows/account/list/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"search":         strings.Join(accountNames, ","),
			"fields":         response.WGAccountList{}.Field(),
			"type":           "exact",
		},
	)

	return res.Data, err
}

func (w *Wargaming) AccountListForSearch(prefix string) (domain.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w.config.URL+"/wows/account/list/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"search":         prefix,
			"fields":         response.WGAccountList{}.Field(),
			"limit":          "10",
		},
	)

	return res.Data, err
}

func (w *Wargaming) ClansAccountInfo(accountIDs []int) (domain.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[response.WGClansAccountInfo](
		w.config.URL+"/wows/clans/accountinfo/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         response.WGClansAccountInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) ClansInfo(clanIDs []int) (domain.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return domain.WGClansInfo{}, nil
	}

	res, err := request[response.WGClansInfo](
		w.config.URL+"/wows/clans/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         response.WGClansInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) ShipsStats(accountID int) (domain.WGShipsStats, error) {
	res, err := request[response.WGShipsStats](
		w.config.URL+"/wows/ships/stats/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"account_id":     strconv.Itoa(accountID),
			"fields":         response.WGShipsStats{}.Field(),
			"extra":          "pvp_solo,pvp_div2,pvp_div3",
		},
	)

	return res.Data, err
}

func (w *Wargaming) EncycShips(pageNo int) (domain.WGEncycShips, int, error) {
	res, err := request[response.WGEncycShips](
		w.config.URL+"/wows/encyclopedia/ships/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"fields":         response.WGEncycShips{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
	)

	return res.Data, res.Meta.PageTotal, err
}

func (w *Wargaming) EncycInfo() (domain.WGEncycInfoData, error) {
	res, err := request[response.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"fields":         response.WGEncycInfo{}.Field(),
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleArenas() (domain.WGBattleArenas, error) {
	res, err := request[response.WGBattleArenas](
		w.config.URL+"/wows/encyclopedia/battlearenas/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"fields":         response.WGBattleArenas{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *Wargaming) BattleTypes() (domain.WGBattleTypes, error) {
	res, err := request[response.WGBattleTypes](
		w.config.URL+"/wows/encyclopedia/battletypes/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": w.appid,
			"fields":         response.WGBattleTypes{}.Field(),
			"language":       "ja",
		},
	)

	return res.Data, err
}

func (w *Wargaming) Test(appid string) (bool, error) {
	_, err := request[response.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		w.config.Retry,
		w.config.Timeout,
		map[string]string{
			"application_id": appid,
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
	operation := func() (webapi.Response[T], error) {
		res, err := webapi.GetRequest[T](rawURL, timeout, queries)
		if err != nil {
			return res, err
		}

		if res.Body.GetStatus() == "error" {
			// Note:
			// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
			message := res.Body.GetError().Message
			if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
				return res, failure.New(apperr.WGAPITemporaryUnavaillalble)
			}

			return res, backoff.Permanent(failure.New(apperr.WGAPIError))
		}

		return res, nil
	}
	res, err := backoff.RetryWithData(operation, b)

	return res.Body, failure.Wrap(err, apperr.ToRequestErrorContext(res))
}
