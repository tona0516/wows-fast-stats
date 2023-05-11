package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

type Wargaming struct {
	AppID string
}

func (w *Wargaming) AccountInfo(accountIDs []int) (vo.WGAccountInfo, error) {
	accountIDsString := make([]string, 0)
	for i := range accountIDs {
		accountIDsString = append(accountIDsString, strconv.Itoa(accountIDs[i]))
	}
	u := buildURL(
		"/wows/account/info/",
		map[string]string{
			"application_id": w.AppID,
			"account_id":     strings.Join(accountIDsString, ","),
			"fields":         vo.WGAccountInfoData{}.Field(),
		},
	)

	return request[vo.WGAccountInfo](u, apperr.Wg.AccountInfo)
}

func (w *Wargaming) AccountList(accountNames []string) (vo.WGAccountList, error) {
	u := buildURL(
		"/wows/account/list/",
		map[string]string{
			"application_id": w.AppID,
			"search":         strings.Join(accountNames, ","),
			"fields":         vo.WGAccountListData{}.Field(),
			"type":           "exact",
		},
	)

	return request[vo.WGAccountList](u, apperr.Wg.AccountList)
}

func (w *Wargaming) ClansAccountInfo(accountIDs []int) (vo.WGClansAccountInfo, error) {
	accountIDsString := make([]string, 0)
	for i := range accountIDs {
		accountIDsString = append(accountIDsString, strconv.Itoa(accountIDs[i]))
	}

	u := buildURL(
		"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": w.AppID,
			"account_id":     strings.Join(accountIDsString, ","),
			"fields":         vo.WGClansAccountInfoData{}.Field(),
		},
	)

	return request[vo.WGClansAccountInfo](u, apperr.Wg.ClansAccountInfo)
}

func (w *Wargaming) ClansInfo(clanIDs []int) (vo.WGClansInfo, error) {
	clanIDsString := make([]string, 0)
	for i := range clanIDs {
		clanIDsString = append(clanIDsString, strconv.Itoa(clanIDs[i]))
	}

	u := buildURL(
		"/wows/clans/info/",
		map[string]string{
			"application_id": w.AppID,
			"clan_id":        strings.Join(clanIDsString, ","),
			"fields":         vo.WGClansInfoData{}.Field(),
		},
	)

	return request[vo.WGClansInfo](u, apperr.Wg.ClansInfo)
}

func (w *Wargaming) EncyclopediaShips(pageNo int) (vo.WGEncyclopediaShips, error) {
	u := buildURL(
		"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGEncyclopediaShipsData{}.Field(),
			"language":       "en",
			"page_no":        strconv.Itoa(pageNo),
		},
	)

	return request[vo.WGEncyclopediaShips](u, apperr.Wg.EncyclopediaShips)
}

func (w *Wargaming) ShipsStats(accountID int) (vo.WGShipsStats, error) {
	u := buildURL(
		"/wows/ships/stats/",
		map[string]string{
			"application_id": w.AppID,
			"account_id":     strconv.Itoa(accountID),
			"fields":         vo.WGShipsStatsData{}.Field(),
		},
	)

	return request[vo.WGShipsStats](u, apperr.Wg.ShipsStats)
}

func (w *Wargaming) EncyclopediaInfo() (vo.WGEncyclopediaInfo, error) {
	u := buildURL(
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGEncyclopediaInfoData{}.Field(),
		},
	)

	return request[vo.WGEncyclopediaInfo](u, apperr.Wg.EncyclopediaInfo)
}

func (w *Wargaming) BattleArenas() (vo.WGBattleArenas, error) {
	u := buildURL(
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGBattleArenasData{}.Field(),
			"language":       "ja",
		},
	)

	return request[vo.WGBattleArenas](u, apperr.Wg.BattleArenas)
}

func (w *Wargaming) BattleTypes() (vo.WGBattleTypes, error) {
	u := buildURL(
		"/wows/encyclopedia/battletypes/",
		map[string]string{
			"application_id": w.AppID,
			"fields":         vo.WGBattleTypesData{}.Field(),
			"language":       "ja",
		},
	)

	return request[vo.WGBattleTypes](u, apperr.Wg.BattleTypes)
}

func buildURL(path string, query map[string]string) *url.URL {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.worldofwarships.asia"
	u.Path = path
	q := u.Query()
	for key, value := range query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u
}

func request[T vo.WGResponse](u *url.URL, errDetail apperr.AppError) (T, error) {
	client := APIClient[T]{}

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
	operation := func() (T, error) {
		res, err := client.GetRequest(u.String())

		// Note:
		// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
		message := res.GetError().Message
		if message == "REQUEST_LIMIT_EXCEEDED" || message == "SOURCE_NOT_AVAILABLE" {
			//nolint:goerr113
			return res, errDetail.WithRaw(fmt.Errorf(message))
		}

		return res, backoff.Permanent(err)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res, errors.WithStack(errDetail.WithRaw(err))
	}
	if res.GetStatus() == "error" {
		//nolint:goerr113
		return res, errors.WithStack(errDetail.WithRaw(fmt.Errorf(res.GetError().Message)))
	}

	return res, err
}
