package infra

import (
	"context"
	"wfs/backend/config"
	"wfs/backend/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type ClanClient struct {
	client *req.Client
}

func NewClanClient(i do.Injector) (*ClanClient, error) {
	config := do.MustInvoke[config.Config](i)
	return &ClanClient{
		client: req.C().
			SetBaseURL(config.ClanClient.URL).
			SetCommonRetryCount(config.ClanClient.RetryCount).
			SetTimeout(config.ClanClient.Timeout),
	}, nil
}

func (c *ClanClient) ClanAutoComplete(ctx context.Context, search string) (data.ClanAutocomplete, error) {
	var result data.ClanAutocomplete
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		SetQueryParams(map[string]string{
			"search": search,
			"type":   "clans",
		}).
		Get("/api/search/autocomplete/")
	if err != nil {
		return result, failure.Translate(err, data.ErrClanAPI)
	}

	if resp.IsErrorState() {
		return result, failure.New(data.ErrClanAPI, failure.Context{
			"status_code": resp.Status,
			"body":        resp.String(),
		})
	}

	return result, nil
}
