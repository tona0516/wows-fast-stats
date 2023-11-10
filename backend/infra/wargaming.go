package infra

import (
	"strconv"
	"strings"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/infra/response"

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
		vo.NewPair("application_id", w.appid),
		vo.NewPair("account_id", strings.Join(strAccountIDs, ",")),
		vo.NewPair("fields", response.WGAccountInfo{}.Field()),
		vo.NewPair("extra", "statistics.pvp_solo"),
	)

	return res.Data, failure.Wrap(err)
}

func (w *Wargaming) AccountList(accountNames []string) (domain.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w.config.URL+"/wows/account/list/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", w.appid),
		vo.NewPair("search", strings.Join(accountNames, ",")),
		vo.NewPair("fields", response.WGAccountList{}.Field()),
		vo.NewPair("type", "exact"),
	)

	return res.Data, failure.Wrap(err)
}

func (w *Wargaming) AccountListForSearch(prefix string) (domain.WGAccountList, error) {
	res, err := request[response.WGAccountList](
		w.config.URL+"/wows/account/list/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", w.appid),
		vo.NewPair("search", prefix),
		vo.NewPair("fields", response.WGAccountList{}.Field()),
		vo.NewPair("limit", "10"),
	)

	return res.Data, failure.Wrap(err)
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
		vo.NewPair("application_id", w.appid),
		vo.NewPair("account_id", strings.Join(strAccountIDs, ",")),
		vo.NewPair("fields", response.WGClansAccountInfo{}.Field()),
	)

	return res.Data, failure.Wrap(err)
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
		vo.NewPair("application_id", w.appid),
		vo.NewPair("clan_id", strings.Join(strClanIDs, ",")),
		vo.NewPair("fields", response.WGClansInfo{}.Field()),
	)

	return res.Data, failure.Wrap(err)
}

func (w *Wargaming) ShipsStats(accountID int) (domain.WGShipsStats, error) {
	res, err := request[response.WGShipsStats](
		w.config.URL+"/wows/ships/stats/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", w.appid),
		vo.NewPair("account_id", strconv.Itoa(accountID)),
		vo.NewPair("fields", response.WGShipsStats{}.Field()),
		vo.NewPair("extra", "pvp_solo"),
	)

	return res.Data, failure.Wrap(err)
}

func (w *Wargaming) EncycShips(pageNo int) (domain.WGEncycShips, int, error) {
	res, err := request[response.WGEncycShips](
		w.config.URL+"/wows/encyclopedia/ships/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", response.WGEncycShips{}.Field()),
		vo.NewPair("language", "ja"),
		vo.NewPair("page_no", strconv.Itoa(pageNo)),
	)

	return res.Data, res.Meta.PageTotal, failure.Wrap(err)
}

func (w *Wargaming) EncycInfo() (domain.WGEncycInfoData, error) {
	res, err := request[response.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", response.WGEncycInfo{}.Field()),
	)

	return res.Data, failure.Wrap(err)
}

func (w *Wargaming) BattleArenas() (domain.WGBattleArenas, error) {
	res, err := request[response.WGBattleArenas](
		w.config.URL+"/wows/encyclopedia/battlearenas/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", response.WGBattleArenas{}.Field()),
		vo.NewPair("language", "ja"),
	)

	return res.Data, failure.Wrap(err)
}

func (w *Wargaming) BattleTypes() (domain.WGBattleTypes, error) {
	res, err := request[response.WGBattleTypes](
		w.config.URL+"/wows/encyclopedia/battletypes/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", response.WGBattleTypes{}.Field()),
		vo.NewPair("language", "ja"),
	)

	return res.Data, failure.Wrap(err)
}

func (w *Wargaming) Test(appid string) (bool, error) {
	_, err := request[response.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		w.config.Retry,
		w.config.Timeout,
		vo.NewPair("application_id", appid),
		vo.NewPair("fields", response.WGEncycInfo{}.Field()),
	)

	return err == nil, failure.Wrap(err)
}

func request[T response.WGResponse](
	rawURL string,
	retry uint64,
	timeout time.Duration,
	query ...vo.Pair,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retry)
	operation := func() (APIResponse[T], error) {
		res, err := getRequest[T](rawURL, timeout, query...)
		if err != nil {
			return res, failure.Wrap(err)
		}

		if res.Body.GetStatus() == "error" {
			// Note:
			// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
			message := res.Body.GetError().Message
			if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
				return res, failure.New(apperr.WGAPITemporaryUnavaillalble)
			}

			return res, failure.New(apperr.WGAPIError)
		}

		return res, nil
	}
	res, err := backoff.RetryWithData(operation, b)

	errCtx := failure.Context{
		"url":         res.FullURL,
		"status_code": strconv.Itoa(res.StatusCode),
		"body":        string(res.ByteBody),
	}

	return res.Body, failure.Wrap(err, errCtx)
}
