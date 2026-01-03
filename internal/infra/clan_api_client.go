package infra

import (
	"wfs/internal/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ClanApiClient interface {
	ClanAutoComplete(search string) (data.ClanAutocomplete, error)
}

type clanApiClient struct {
	client *req.Client
}

func NewClansApiClient(ac apiConfig) ClanApiClient {
	return &clanApiClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout),
	}
}

func (c *clanApiClient) ClanAutoComplete(search string) (data.ClanAutocomplete, error) {
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
