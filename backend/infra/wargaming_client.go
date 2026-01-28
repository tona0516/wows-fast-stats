package infra

import (
	"context"
	"encoding/json"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"wfs/backend/config"
	"wfs/backend/core"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"go.uber.org/ratelimit"
)

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

				var body core.WGResponseCommon[any]
				if err := json.Unmarshal(resp.Bytes(), &body); err != nil {
					return true
				}

				err = handleError(resp, err, body)
				return failure.Is(err, core.ErrWGAPITemporaryUnavailable)
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
	accountIDs []core.AccountID,
) (core.WGAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(int(v))
	}

	return request[core.WGAccountInfo](
		ctx,
		c,
		"/wows/account/info/",
		map[string]string{
			"account_id": strings.Join(strAccountIDs, ","),
			"fields":     c.fieldQuery(reflect.TypeFor[core.WGAccountInfoData]()),
			"extra": strings.Join([]string{
				"statistics.pvp_solo",
				"statistics.pvp_div2",
				"statistics.pvp_div3",
				"statistics.rank_solo",
			}, ","),
		},
	)
}

func (c *WargamingClient) AccountList(
	ctx context.Context,
	accountNames []string,
) (core.WGAccountList, error) {
	return request[core.WGAccountList](
		ctx,
		c,
		"/wows/account/list/",
		map[string]string{
			"search": strings.Join(accountNames, ","),
			"fields": c.fieldQuery(reflect.TypeFor[core.WGAccountListData]()),
			"type":   "exact",
		},
	)
}

func (c *WargamingClient) ClansAccountInfo(
	ctx context.Context,
	accountIDs []core.AccountID,
) (core.WGClansAccountInfo, error) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(int(v))
	}

	return request[core.WGClansAccountInfo](
		ctx,
		c,
		"/wows/clans/accountinfo/",
		map[string]string{
			"account_id": strings.Join(strAccountIDs, ","),
			"fields":     c.fieldQuery(reflect.TypeFor[core.WGClansAccountInfoData]()),
		},
	)
}

func (c *WargamingClient) ClansInfo(
	ctx context.Context,
	clanIDs []core.ClanID,
) (core.WGClansInfo, error) {
	strClanIDs := make([]string, len(clanIDs))
	for i, v := range clanIDs {
		strClanIDs[i] = strconv.Itoa(int(v))
	}

	if len(strClanIDs) == 0 {
		return core.WGClansInfo{}, nil
	}

	return request[core.WGClansInfo](
		ctx,
		c,
		"/wows/clans/info/",
		map[string]string{
			"clan_id": strings.Join(strClanIDs, ","),
			"fields":  c.fieldQuery(reflect.TypeFor[core.WGClansInfoData]()),
		},
	)
}

func (c *WargamingClient) ShipsStats(
	ctx context.Context,
	accountID core.AccountID,
) (core.WGShipsStats, error) {
	return request[core.WGShipsStats](
		ctx,
		c,
		"/wows/ships/stats/",
		map[string]string{
			"account_id": strconv.Itoa(int(accountID)),
			"fields":     c.fieldQuery(reflect.TypeFor[core.WGShipsStatsData]()),
			"extra": strings.Join([]string{
				"pvp_solo",
				"pvp_div2",
				"pvp_div3",
				"rank_solo",
			}, ","),
		},
	)
}

func (c *WargamingClient) EncycShips(
	ctx context.Context,
	pageNo int,
) (core.WGEncycShips, error) {
	return request[core.WGEncycShips](
		ctx,
		c,
		"/wows/encyclopedia/ships/",
		map[string]string{
			"fields":   c.fieldQuery(reflect.TypeFor[core.WGEncycShipsData]()),
			"language": "ja",
			"page_no":  strconv.Itoa(pageNo),
		},
	)
}

func (c *WargamingClient) ShipsBadges(
	ctx context.Context,
	accountID core.AccountID,
) (core.WGShipsBadges, error) {
	return request[core.WGShipsBadges](
		ctx,
		c,
		"/wows/ships/badges/",
		map[string]string{
			"account_id": strconv.Itoa(int(accountID)),
			"fields":     c.fieldQuery(reflect.TypeFor[core.WGShipsBadgesData]()),
		},
	)
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

func request[T core.WGResponse](
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
	err = handleError(resp, err, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func handleError[T core.WGResponse](resp *req.Response, err error, result T) error {
	if err != nil {
		return failure.Translate(err, core.ErrWGAPI)
	}

	if resp.IsErrorState() {
		return failure.New(core.ErrWGAPI, failure.Context{
			"status_code": resp.Status,
			"body":        string(resp.Bytes()),
		})
	}

	errResp := result.GetError()
	errCtx := failure.Context{
		"status":        result.GetStatus(),
		"error.code":    strconv.Itoa(errResp.Code),
		"error.message": errResp.Message,
		"error.field":   errResp.Field,
		"error.value":   errResp.Value,
	}

	if result.GetStatus() == "error" {
		// Note:
		// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
		if slices.Contains(temporaryUnavaillalbleMessages, errResp.Message) {
			return failure.New(core.ErrWGAPITemporaryUnavailable, errCtx)
		}

		return failure.New(core.ErrWGAPI, errCtx)
	}

	return nil
}
