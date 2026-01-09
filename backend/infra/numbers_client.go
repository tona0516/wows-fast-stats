package infra

import (
	"encoding/json"
	"wfs/backend/config"
	"wfs/backend/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type NumbersClient struct {
	client *req.Client
}

func NewNumbersClient(i do.Injector) (*NumbersClient, error) {
	config := do.MustInvoke[config.Config](i)
	return &NumbersClient{
		client: req.C().
			SetBaseURL(config.NumbersClient.URL).
			SetCommonRetryCount(config.NumbersClient.RetryCount).
			SetTimeout(config.NumbersClient.Timeout).
			EnableInsecureSkipVerify(),
	}, nil
}

func (c *NumbersClient) ExpectedStats() (data.NSExpectedStats, error) {
	var result data.NSExpectedStats

	resp, err := c.client.R().
		Get("/personal/rating/expected/json/")
	if err != nil {
		return result, failure.Wrap(err)
	}

	if resp.IsErrorState() {
		return result, failure.Wrap(ErrErrorResponse)
	}

	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return result, failure.Wrap(err)
	}

	return result, nil
}
