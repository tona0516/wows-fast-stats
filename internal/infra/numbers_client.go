package infra

import (
	"encoding/json"
	"wfs/internal/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type NumbersClient struct {
	client *req.Client
}

func NewNumbersClient(ac apiConfig) *NumbersClient {
	return &NumbersClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout).
			EnableInsecureSkipVerify(),
	}
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
