package repo

import (
	"changeme/backend/vo"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func buildUrl(path string, query map[string]string) *url.URL {
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

type Wargaming struct {
	AppID string
}

func (w *Wargaming) GetAccountInfo(accountIDs []int) (vo.WGAccountInfo, error) {
	accountIDsString := make([]string, 0)
	for i := range accountIDs {
		accountIDsString = append(accountIDsString, strconv.Itoa(accountIDs[i]))
	}
	u := buildUrl(
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

	client := ApiClient[vo.WGAccountInfo]{}
	res, err := client.GetRequest(u.String())
	if res.Status == "error" {
		return res, errors.New(res.Error.Message)
	}
	return res, err
}

func (w *Wargaming) GetAccountList(accountNames []string) (vo.WGAccountList, error) {
	u := buildUrl(
		"/wows/account/list/",
		map[string]string{
			"application_id": w.AppID,
			"search":         strings.Join(accountNames, ","),
			"fields":         strings.Join([]string{"account_id", "nickname"}, ","),
			"type":           "exact",
		},
	)

	client := ApiClient[vo.WGAccountList]{}
	res, err := client.GetRequest(u.String())
	if res.Status == "error" {
		return res, errors.New(res.Error.Message)
	}
	return res, err
}

func (w *Wargaming) GetClansAccountInfo(accountIDs []int) (vo.WGClansAccountInfo, error) {
	accountIDsString := make([]string, 0)
	for i := range accountIDs {
		accountIDsString = append(accountIDsString, strconv.Itoa(accountIDs[i]))
	}

	u := buildUrl(
		"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": w.AppID,
			"account_id":     strings.Join(accountIDsString, ","),
			"fields":         "clan_id",
		},
	)

	client := ApiClient[vo.WGClansAccountInfo]{}
	res, err := client.GetRequest(u.String())
	if res.Status == "error" {
		return res, errors.New(res.Error.Message)
	}
	return res, err
}

func (w *Wargaming) GetClansInfo(clanIDs []int) (vo.WGClansInfo, error) {
	clanIDsString := make([]string, 0)
	for i := range clanIDs {
		clanIDsString = append(clanIDsString, strconv.Itoa(clanIDs[i]))
	}

	u := buildUrl(
		"/wows/clans/info/",
		map[string]string{
			"application_id": w.AppID,
			"clan_id":        strings.Join(clanIDsString, ","),
			"fields":         "tag",
		},
	)

	client := ApiClient[vo.WGClansInfo]{}
	res, err := client.GetRequest(u.String())
	if res.Status == "error" {
		return res, errors.New(res.Error.Message)
	}
	return res, err
}

func (w *Wargaming) GetEncyclopediaShips(pageNo int) (vo.WGEncyclopediaShips, error) {
	u := buildUrl(
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

	client := ApiClient[vo.WGEncyclopediaShips]{}
	res, err := client.GetRequest(u.String())
	if res.Status == "error" {
		return res, errors.New(res.Error.Message)
	}
	return res, err
}

func (w *Wargaming) GetShipsStats(accountID int) (vo.WGShipsStats, error) {
	u := buildUrl(
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

	client := ApiClient[vo.WGShipsStats]{}
	res, err := client.GetRequest(u.String())
	if res.Status == "error" {
		return res, errors.New(res.Error.Message)
	}
	return res, err
}

func (w *Wargaming) GetEncyclopediaInfo() (vo.WGEncyclopediaInfo, error) {
	u := buildUrl(
		"/wows/encyclopedia/info/",
		map[string]string{
			"application_id": w.AppID,
			"fields":         "game_version",
		},
	)

	client := ApiClient[vo.WGEncyclopediaInfo]{}
	res, err := client.GetRequest(u.String())
	if res.Status == "error" {
		return res, errors.New(res.Error.Message)
	}
	return res, err
}
