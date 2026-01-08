package infra

import (
	"wfs/internal/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type ClanClient struct {
	client *req.Client
}

func NewClanClient(ac apiConfig) *ClanClient {
	return &ClanClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout),
	}
}

func (c *ClanClient) ClanAutoComplete(search string) (data.ClanAutocomplete, error) {
	var result data.ClanAutocomplete
	resp, err := c.client.R().
		SetSuccessResult(&result).
		SetQueryParams(map[string]string{
			"search": search,
			"type":   "clans",
		}).
		Get("/api/search/autocomplete/")
	if err != nil {
		return result, failure.Wrap(err)
	}

	if resp.IsErrorState() {
		return result, failure.Wrap(ErrErrorResponse)
	}

	return result, nil
}
