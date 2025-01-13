package infra

import (
	"encoding/json"
	"wfs/backend/apperr"
	"wfs/backend/infra/response"
	"wfs/backend/infra/webapi"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
)

type Numbers struct {
	config RequestConfig
}

func NewNumbers(config RequestConfig) *Numbers {
	return &Numbers{config: config}
}

func (n *Numbers) expectedStats() (response.ExpectedStats, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.config.Retry)
	operation := func() (response.NSExpectedStats, error) {
		var result response.NSExpectedStats

		_, body, err := webapi.NewClient(n.config.URL,
			webapi.WithPath("/personal/rating/expected/json/"),
			webapi.WithTimeout(n.config.Timeout),
			webapi.WithIsInsecure(true), // workaround for expired SSL certificate
		).GET()
		if err != nil {
			return result, err
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return result, failure.Wrap(err)
		}

		return result, nil
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res.Data, failure.Translate(err, apperr.NumbersAPIError)
	}

	return res.Data, failure.Wrap(err)
}
