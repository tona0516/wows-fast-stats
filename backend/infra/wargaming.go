package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"fmt"
	"net/url"
	"strconv"
	"strings"

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
			"fields": strings.Join([]string{
				"hidden_profile",
				"statistics.pvp.battles",
				"statistics.pvp.wins",
				"statistics.pvp.frags",
				"statistics.pvp.damage_dealt",
				"statistics.pvp.xp",
				"statistics.pvp.survived_battles",
				"statistics.pvp.survived_wins",
			}, ","),
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
			"fields":         strings.Join([]string{"account_id", "nickname"}, ","),
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
			"fields":         "clan_id",
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
			"fields":         "tag",
		},
	)

	return request[vo.WGClansInfo](u, apperr.Wg.ClansInfo)
}

func (w *Wargaming) EncyclopediaShips(pageNo int) (vo.WGEncyclopediaShips, error) {
	u := buildURL(
		"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": w.AppID,
			"fields": strings.Join([]string{
				"name",
				"tier",
				"type",
				"nation",
			}, ","),
			"language": "en",
			"page_no":  strconv.Itoa(pageNo),
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
			"fields": strings.Join([]string{
				"ship_id",
				"pvp.battles",
				"pvp.wins",
				"pvp.frags",
				"pvp.damage_dealt",
				"pvp.xp",
				"pvp.survived_battles",
				"pvp.survived_wins",
			}, ","),
		},
	)

	return request[vo.WGShipsStats](u, apperr.Wg.ShipsStats)
}

func (w *Wargaming) EncyclopediaInfo() (vo.WGEncyclopediaInfo, error) {
	u := buildURL(
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": w.AppID,
			"fields":         "game_version",
		},
	)

	return request[vo.WGEncyclopediaInfo](u, apperr.Wg.EncyclopediaInfo)
}

func (w *Wargaming) BattleArenas() (vo.WGBattleArenas, error) {
	u := buildURL(
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"application_id": w.AppID,
			"fields":         "name",
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
			"fields":         "name",
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
	res, err := client.GetRequest(u.String())
	if err != nil {
		return res, errors.WithStack(errDetail.WithRaw(err))
	}
	if res.GetStatus() == "error" {
		//nolint:goerr113
		return res, errors.WithStack(errDetail.WithRaw(fmt.Errorf(res.GetError().Message)))
	}

	return res, err
}
