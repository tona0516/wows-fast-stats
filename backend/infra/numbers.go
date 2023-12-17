//nolint:dupl
package infra

import (
	"net/http"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/domain"
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

func (n *Numbers) ExpectedStats() (domain.NSExpectedStats, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.config.Retry)
	operation := func() (webapi.Response[domain.NSExpectedStats], error) {
		return webapi.GetRequest[domain.NSExpectedStats](n.config.URL+"/personal/rating/expected/json/", n.config.Timeout)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res.Body, failure.Wrap(err)
	}

	if res.StatusCode != http.StatusOK {
		return res.Body, failure.New(
			apperr.NumbersAPIFetchExpectedStatsError,
			failure.Context{
				"status_code": strconv.Itoa(res.StatusCode),
				"body":        string(res.ByteBody),
			},
		)
	}

	return res.Body, nil
}
