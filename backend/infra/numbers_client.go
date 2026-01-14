package infra

import (
	"context"
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

func (c *NumbersClient) ExpectedStats(ctx context.Context) (data.NSExpectedStats, error) {
	var result data.NSExpectedStats

	resp, err := c.client.R().
		SetContext(ctx).
		Get("/personal/rating/expected/json/")
	if err != nil {
		return result, failure.Translate(err, data.ErrNumbersAPI)
	}

	if resp.IsErrorState() {
		return result, failure.New(data.ErrNumbersAPI, failure.Context{
			"status_code": resp.Status,
			"body":        resp.String(),
		})
	}

	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return result, failure.Translate(err, data.ErrNumbersAPI)
	}

	return result, nil
}
