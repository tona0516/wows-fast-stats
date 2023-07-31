package infra

import (
	"strconv"
	"strings"
	"wfs/backend/apperr"

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
		w.config.URL+"/wows/account/info/",
		map[string]string{
			"application_id": w.appid,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         domain.WGAccountInfoData{}.Field(),
			"extra":          "statistics.pvp_solo",
		},
		w.config.Retry,
	)
}

func (w *Wargaming) AccountList(accountNames []string) (domain.WGAccountList, error) {
	return request[domain.WGAccountList](
		w.config.URL+"/wows/account/list/",
		map[string]string{
			"application_id": w.appid,
			"search":         strings.Join(accountNames, ","),
			"fields":         domain.WGAccountListData{}.Field(),
			"type":           "exact",
		},
		w.config.Retry,
	)
}

func (w *Wargaming) AccountListForSearch(prefix string) (domain.WGAccountList, error) {
	return request[domain.WGAccountList](
		w.config.URL+"/wows/account/list/",
		map[string]string{
			"application_id": w.appid,
			"search":         prefix,
			"fields":         domain.WGAccountListData{}.Field(),
		},
		w.config.Retry,
	)
}

func (w *Wargaming) ClansAccountInfo(accountIDs []int) (domain.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	return request[domain.WGClansAccountInfo](
		w.config.URL+"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": w.appid,
			"account_id":     strings.Join(strAccountIDs, ","),
			"fields":         domain.WGClansAccountInfoData{}.Field(),
		},
		w.config.Retry,
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
		w.config.URL+"/wows/clans/info/",
		map[string]string{
			"application_id": w.appid,
			"clan_id":        strings.Join(strClanIDs, ","),
			"fields":         domain.WGClansInfoData{}.Field(),
		},
		w.config.Retry,
	)
}

func (w *Wargaming) ShipsStats(accountID int) (domain.WGShipsStats, error) {
	return request[domain.WGShipsStats](
		w.config.URL+"/wows/ships/stats/",
		map[string]string{
			"application_id": w.appid,
			"account_id":     strconv.Itoa(accountID),
			"fields":         domain.WGShipsStatsData{}.Field(),
			"extra":          "pvp_solo",
		},
		w.config.Retry,
	)
}

func (w *Wargaming) EncycShips(pageNo int) (domain.WGEncycShips, error) {
	return request[domain.WGEncycShips](
		w.config.URL+"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": w.appid,
			"fields":         domain.WGEncyclopediaShipsData{}.Field(),
			"language":       "ja",
			"page_no":        strconv.Itoa(pageNo),
		},
		w.config.Retry,
	)
}

func (w *Wargaming) EncycInfo() (domain.WGEncycInfo, error) {
	return request[domain.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": w.appid,
			"fields":         domain.WGEncyclopediaInfoData{}.Field(),
		},
		w.config.Retry,
	)
}

func (w *Wargaming) BattleArenas() (domain.WGBattleArenas, error) {
	return request[domain.WGBattleArenas](
		w.config.URL+"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": w.appid,
			"fields":         domain.WGBattleArenasData{}.Field(),
			"language":       "ja",
		},
		w.config.Retry,
	)
}

func (w *Wargaming) BattleTypes() (domain.WGBattleTypes, error) {
	return request[domain.WGBattleTypes](
		w.config.URL+"/wows/encyclopedia/battletypes/",
		map[string]string{
			"application_id": w.appid,
			"fields":         domain.WGBattleTypesData{}.Field(),
			"language":       "ja",
		},
		w.config.Retry,
	)
}

func (w *Wargaming) Test(appid string) (bool, error) {
	_, err := request[domain.WGEncycInfo](
		w.config.URL+"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": appid,
			"fields":         domain.WGEncyclopediaInfoData{}.Field(),
		},
		w.config.Retry,
	)

	return err == nil, err
}

func request[T domain.WGResponse](
	rawURL string,
	query map[string]string,
	retry uint64,
) (T, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retry)
	operation := func() (APIResponse[T], error) {
		res, err := getRequest[T](rawURL, query, retry)

		// Note:
		// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
		message := res.Body.GetError().Message
		if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
			return res, apperr.New(apperr.ErrWGAPITemporaryUnavaillalble, errors.New(res.BodyString))
		}

		return res, backoff.Permanent(err)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res.Body, err
	}

	if res.Body.GetStatus() == "error" {
		return res.Body, apperr.New(apperr.ErrWGAPI, errors.New(res.BodyString))
	}

	return res.Body, err
}
