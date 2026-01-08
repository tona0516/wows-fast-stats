package infra

import (
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"strings"
	"wfs/internal/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"go.uber.org/ratelimit"
)

//nolint:gochecknoglobals
var temporaryUnavaillalbleMessages = []string{
	"REQUEST_LIMIT_EXCEEDED",
	"SOURCE_NOT_AVAILABLE",
}

type WargamingClient struct {
	client *req.Client
}

func NewWargamingClient(
	appID string,
	ac apiConfig,
	limiter ratelimit.Limiter,
) *WargamingClient {
	return &WargamingClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout).
			SetCommonQueryParam("application_id", appID).
			SetCommonRetryCondition(func(resp *req.Response, err error) bool {
				if err != nil {
					return true
				}

				var body data.WGResponseCommon[any]
				if err := json.Unmarshal(resp.Bytes(), &body); err == nil {
					err := convertError(body.Status, body.Error.Message)
					if errors.Is(err, ErrTemporaryUnavaillalble) {
						return true
					}
				}

				return false
			}).
			OnBeforeRequest(func(client *req.Client, req *req.Request) error {
				limiter.Take()
				return nil
			}).
			SetCommonRetryHook(func(resp *req.Response, err error) {
				limiter.Take()
			}),
	}
}

func (c *WargamingClient) AccountInfo(accountIDs []int) (data.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[data.WGAccountInfo](
		c,
		"/wows/account/info/",
		map[string]string{
			"account_id": strings.Join(strAccountIDs, ","),
			"fields":     data.WGAccountInfo{}.Field(),
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

func (c *WargamingClient) AccountList(accountNames []string) (data.WGAccountList, error) {
	res, err := request[data.WGAccountList](
		c,
		"/wows/account/list/",
		map[string]string{
			"search": strings.Join(accountNames, ","),
			"fields": data.WGAccountList{}.Field(),
			"type":   "exact",
		},
	)

	return res, err
}

func (c *WargamingClient) ClansAccountInfo(accountIDs []int) (data.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	res, err := request[data.WGClansAccountInfo](
		c,
		"/wows/clans/accountinfo/",
		map[string]string{
			"account_id": strings.Join(strAccountIDs, ","),
			"fields":     data.WGClansAccountInfo{}.Field(),
		},
	)

	return res, err
}

func (c *WargamingClient) ClansInfo(clanIDs []int) (data.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(v)
	}

	if len(strClanIDs) == 0 {
		return data.WGClansInfo{}, nil
	}

	res, err := request[data.WGClansInfo](
		c,
		"/wows/clans/info/",
		map[string]string{
			"clan_id": strings.Join(strClanIDs, ","),
			"fields":  data.WGClansInfo{}.Field(),
		},
	)

	return res, err
}

func (c *WargamingClient) ShipsStats(accountID int) (data.WGShipsStats, error) {
	res, err := request[data.WGShipsStats](
		c,
		"/wows/ships/stats/",
		map[string]string{
			"account_id": strconv.Itoa(accountID),
			"fields":     data.WGShipsStats{}.Field(),
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

func (c *WargamingClient) EncycShips(pageNo int) (data.WGEncycShips, error) {
	res, err := request[data.WGEncycShips](
		c,
		"/wows/encyclopedia/ships/",
		map[string]string{
			"fields":   data.WGEncycShips{}.Field(),
			"language": "ja",
			"page_no":  strconv.Itoa(pageNo),
		},
	)

	return res, err
}

func (c *WargamingClient) BattleArenas() (data.WGBattleArenas, error) {
	res, err := request[data.WGBattleArenas](
		c,
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"fields":   data.WGBattleArenas{}.Field(),
			"language": "ja",
		},
	)

	return res, err
}

func (c *WargamingClient) BattleTypes() (data.WGBattleTypes, error) {
	res, err := request[data.WGBattleTypes](
		c,
		"/wows/encyclopedia/battletypes/",
		map[string]string{
			"fields":   data.WGBattleTypes{}.Field(),
			"language": "ja",
		},
	)

	return res, err
}

func (c *WargamingClient) ShipsBadges(accountID int) (data.WGShipsBadges, error) {
	res, err := request[data.WGShipsBadges](
		c,
		"/wows/ships/badges/",
		map[string]string{
			"account_id": strconv.Itoa(accountID),
			"fields":     data.WGShipsBadges{}.Field(),
		},
	)

	return res, err
}

func request[T data.WGResponse](
	c *WargamingClient,
	path string,
	queries map[string]string,
) (T, error) {
	var result T

	client := c.client.R().SetSuccessResult(&result)
	for k, v := range queries {
		client.AddQueryParam(k, v)
	}

	resp, err := client.Get(path)
	if err != nil {
		return result, failure.Wrap(err)
	}

	if resp.IsErrorState() {
		return result, failure.Wrap(ErrErrorResponse)
	}

	if err := convertError(result.GetStatus(), result.GetError().Message); err != nil {
		return result, failure.Wrap(err)
	}

	return result, nil
}

func convertError(status string, message string) error {
	if status == "ok" {
		return nil
	}

	// Note:
	// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
	if slices.Contains(temporaryUnavaillalbleMessages, message) {
		return failure.Wrap(ErrTemporaryUnavaillalble)
	}

	return failure.Wrap(ErrErrorResponse)
}
