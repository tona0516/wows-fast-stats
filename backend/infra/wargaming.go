package infra

import (
	"strconv"
	"strings"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/logger"

	"wfs/backend/domain"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
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

	return request[domain.WGAccountInfo](
		w.config.Retry,
		w.config.URL+"/wows/account/info/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("account_id", strings.Join(strAccountIDs, ",")),
		vo.NewPair("fields", domain.WGAccountInfoData{}.Field()),
		vo.NewPair("extra", "statistics.pvp_solo"),
	)
}

func (w *Wargaming) AccountList(accountNames []string) (domain.WGAccountList, error) {
	return request[domain.WGAccountList](
		w.config.Retry,
		w.config.URL+"/wows/account/list/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("search", strings.Join(accountNames, ",")),
		vo.NewPair("fields", domain.WGAccountListData{}.Field()),
		vo.NewPair("type", "exact"),
	)
}

func (w *Wargaming) AccountListForSearch(prefix string) (domain.WGAccountList, error) {
	return request[domain.WGAccountList](
		w.config.Retry,
		w.config.URL+"/wows/account/list/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("search", prefix),
		vo.NewPair("fields", domain.WGAccountListData{}.Field()),
	)
}

func (w *Wargaming) ClansAccountInfo(accountIDs []int) (domain.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	return request[domain.WGClansAccountInfo](
		w.config.Retry,
		w.config.URL+"/wows/clans/accountinfo/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("account_id", strings.Join(strAccountIDs, ",")),
		vo.NewPair("fields", domain.WGClansAccountInfoData{}.Field()),
	)
}

func (w *Wargaming) ClansInfo(clanIDs []int) (domain.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return domain.WGClansInfo{}, nil
	}

	return request[domain.WGClansInfo](
		w.config.Retry,
		w.config.URL+"/wows/clans/info/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("clan_id", strings.Join(strClanIDs, ",")),
		vo.NewPair("fields", domain.WGClansInfoData{}.Field()),
	)
}

func (w *Wargaming) ShipsStats(accountID int) (domain.WGShipsStats, error) {
	return request[domain.WGShipsStats](
		w.config.Retry,
		w.config.URL+"/wows/ships/stats/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("account_id", strconv.Itoa(accountID)),
		vo.NewPair("fields", domain.WGShipsStatsData{}.Field()),
		vo.NewPair("extra", "pvp_solo"),
	)
}

func (w *Wargaming) EncycShips(pageNo int) (domain.WGEncycShips, error) {
	return request[domain.WGEncycShips](
		w.config.Retry,
		w.config.URL+"/wows/encyclopedia/ships/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", domain.WGEncyclopediaShipsData{}.Field()),
		vo.NewPair("language", "ja"),
		vo.NewPair("page_no", strconv.Itoa(pageNo)),
	)
}

func (w *Wargaming) EncycInfo() (domain.WGEncycInfo, error) {
	return request[domain.WGEncycInfo](
		w.config.Retry,
		w.config.URL+"/wows/encyclopedia/info/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", domain.WGEncyclopediaInfoData{}.Field()),
	)
}

func (w *Wargaming) BattleArenas() (domain.WGBattleArenas, error) {
	return request[domain.WGBattleArenas](
		w.config.Retry,
		w.config.URL+"/wows/encyclopedia/battlearenas/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", domain.WGBattleArenasData{}.Field()),
		vo.NewPair("language", "ja"),
	)
}

func (w *Wargaming) BattleTypes() (domain.WGBattleTypes, error) {
	return request[domain.WGBattleTypes](
		w.config.Retry,
		w.config.URL+"/wows/encyclopedia/battletypes/",
		vo.NewPair("application_id", w.appid),
		vo.NewPair("fields", domain.WGBattleTypesData{}.Field()),
		vo.NewPair("language", "ja"),
	)
}

func (w *Wargaming) Test(appid string) (bool, error) {
	_, err := request[domain.WGEncycInfo](
		w.config.Retry,
		w.config.URL+"/wows/encyclopedia/info/",
		vo.NewPair("application_id", appid),
		vo.NewPair("fields", domain.WGEncyclopediaInfoData{}.Field()),
	)

	return err == nil, err
}

func request[T domain.WGResponse](
	retry uint64,
	rawURL string,
	query ...vo.Pair,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retry)
	operation := func() (APIResponse[T], error) {
		res, err := getRequest[T](rawURL, query...)

		if err != nil {
			return res, backoff.Permanent(apperr.New(apperr.ErrWGAPI, err))
		}

		if res.Body.GetStatus() == "error" {
			// Note:
			// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
			message := res.Body.GetError().Message
			if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
				return res, apperr.New(apperr.ErrWGAPITemporaryUnavaillalble, errors.New(message))
			}

			return res, apperr.New(apperr.ErrWGAPI, errors.New(message))
		}

		return res, nil
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil && !errors.Is(err, apperr.ErrWGAPITemporaryUnavaillalble) {
		logger.Error(
			err,
			vo.NewPair("url", rawURL),
			vo.NewPair("status_code", strconv.Itoa(res.StatusCode)),
			vo.NewPair("response_body", string(res.BodyByte)),
		)
	}

	return res.Body, err
}
