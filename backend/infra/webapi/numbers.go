package webapi

import (
	"encoding/json"
	"wfs/backend/apperr"
	"wfs/backend/infra/response"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type Numbers interface {
	ExpectedStats() (response.ExpectedStats, error)
}

type numbers struct {
	config RequestConfig
}

func NewNumbers(config RequestConfig) Numbers {
	return &numbers{config: config}
}

func (n *numbers) ExpectedStats() (response.ExpectedStats, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.config.Retry)
	operation := func() (response.NSExpectedStats, error) {
		var result response.NSExpectedStats

		_, body, err := NewClient(n.config.URL,
			WithPath("/personal/rating/expected/json/"),
			WithTimeout(n.config.Timeout),
			WithIsInsecure(true), // workaround for expired SSL certificate
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
