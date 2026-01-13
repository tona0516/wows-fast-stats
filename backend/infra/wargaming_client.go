package infra

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"wfs/backend/config"
	"wfs/backend/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
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

func NewWargamingClient(i do.Injector) (*WargamingClient, error) {
	config := do.MustInvoke[config.Config](i)
	limiter := ratelimit.New(config.WargamingClient.RateLimitRPS)

	return &WargamingClient{
		client: req.C().
			SetBaseURL(config.WargamingClient.URL).
			SetCommonRetryCount(config.WargamingClient.RetryCount).
			SetTimeout(config.WargamingClient.Timeout).
			SetCommonQueryParam("application_id", config.WargamingClient.AppID).
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
	}, nil
}

func (c *WargamingClient) AccountInfo(
	ctx context.Context,
	accountIDs []data.AccountID,
) (data.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(int(v))
	}

	res, err := request[data.WGAccountInfo](
		ctx,
		c,
		"/wows/account/info/",
		map[string]string{
			"account_id": strings.Join(strAccountIDs, ","),
			"fields":     c.fieldQuery(reflect.TypeFor[data.WGAccountInfoData]()),
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

func (c *WargamingClient) AccountList(
	ctx context.Context,
	accountNames []string,
) (data.WGAccountList, error) {
	res, err := request[data.WGAccountList](
		ctx,
		c,
		"/wows/account/list/",
		map[string]string{
			"search": strings.Join(accountNames, ","),
			"fields": c.fieldQuery(reflect.TypeFor[data.WGAccountListData]()),
			"type":   "exact",
		},
	)

	return res, err
}

func (c *WargamingClient) ClansAccountInfo(
	ctx context.Context,
	accountIDs []data.AccountID,
) (data.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(int(v))
	}

	res, err := request[data.WGClansAccountInfo](
		ctx,
		c,
		"/wows/clans/accountinfo/",
		map[string]string{
			"account_id": strings.Join(strAccountIDs, ","),
			"fields":     c.fieldQuery(reflect.TypeFor[data.WGClansAccountInfoData]()),
		},
	)

	return res, err
}

func (c *WargamingClient) ClansInfo(
	ctx context.Context,
	clanIDs []data.ClanID,
) (data.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(int(v))
	}

	if len(strClanIDs) == 0 {
		return data.WGClansInfo{}, nil
	}

	res, err := request[data.WGClansInfo](
		ctx,
		c,
		"/wows/clans/info/",
		map[string]string{
			"clan_id": strings.Join(strClanIDs, ","),
			"fields":  c.fieldQuery(reflect.TypeFor[data.WGClansInfoData]()),
		},
	)

	return res, err
}

func (c *WargamingClient) ShipsStats(
	ctx context.Context,
	accountID data.AccountID,
) (data.WGShipsStats, error) {
	res, err := request[data.WGShipsStats](
		ctx,
		c,
		"/wows/ships/stats/",
		map[string]string{
			"account_id": strconv.Itoa(int(accountID)),
			"fields":     c.fieldQuery(reflect.TypeFor[data.WGShipsStatsData]()),
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

func (c *WargamingClient) EncycShips(
	ctx context.Context,
	pageNo int,
) (data.WGEncycShips, error) {
	res, err := request[data.WGEncycShips](
		ctx,
		c,
		"/wows/encyclopedia/ships/",
		map[string]string{
			"fields":   c.fieldQuery(reflect.TypeFor[data.WGEncycShipsData]()),
			"language": "ja",
			"page_no":  strconv.Itoa(pageNo),
		},
	)

	return res, err
}

func (c *WargamingClient) BattleArenas(ctx context.Context) (data.WGBattleArenas, error) {
	res, err := request[data.WGBattleArenas](
		ctx,
		c,
		"/wows/encyclopedia/battlearenas/",
		map[string]string{
			"fields":   c.fieldQuery(reflect.TypeFor[data.WGBattleArenasData]()),
			"language": "ja",
		},
	)

	return res, err
}

func (c *WargamingClient) BattleTypes(ctx context.Context) (data.WGBattleTypes, error) {
	res, err := request[data.WGBattleTypes](
		ctx,
		c,
		"/wows/encyclopedia/battletypes/",
		map[string]string{
			"fields":   c.fieldQuery(reflect.TypeFor[data.WGBattleTypesData]()),
			"language": "ja",
		},
	)

	return res, err
}

func (c *WargamingClient) ShipsBadges(
	ctx context.Context,
	accountID data.AccountID,
) (data.WGShipsBadges, error) {
	res, err := request[data.WGShipsBadges](
		ctx,
		c,
		"/wows/ships/badges/",
		map[string]string{
			"account_id": strconv.Itoa(int(accountID)),
			"fields":     c.fieldQuery(reflect.TypeFor[data.WGShipsBadgesData]()),
		},
	)

	return res, err
}

func (c *WargamingClient) fieldQuery(target reflect.Type) string {
	fields := []string{}
	c.fieldQueryRecursive([]string{}, target, &fields)

	return strings.Join(fields, ",")
}

func (c *WargamingClient) fieldQueryRecursive(parentNames []string, t reflect.Type, result *[]string) {
	for i := range t.NumField() {
		field := t.Field(i)

		if field.Type.Kind() == reflect.Struct {
			name := c.toSnakeCase(field.Name)
			c.fieldQueryRecursive(append(parentNames, name), field.Type, result)

			continue
		}

		column := field.Tag.Get("json")
		if len(parentNames) > 0 {
			*result = append(*result, strings.Join(parentNames, ".")+"."+column)
		} else {
			*result = append(*result, column)
		}
	}
}

func (c *WargamingClient) toSnakeCase(s string) string {
	runes := []rune(s)
	result := make([]rune, 0)

	for i, r := range runes {
		if !unicode.IsUpper(r) {
			result = append(result, r)
			continue
		}

		if i > 0 && unicode.IsLower(runes[i-1]) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}

func request[T data.WGResponse](
	ctx context.Context,
	c *WargamingClient,
	path string,
	queries map[string]string,
) (T, error) {
	var result T

	client := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result)

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
