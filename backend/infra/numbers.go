package infra

import (
	"net/http"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
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

func (n *Numbers) ExpectedStats() (model.ExpectedStats, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.config.Retry)
	operation := func() (webapi.Response[any, model.NSExpectedStats], error) {
		res, err := webapi.GetRequest[model.NSExpectedStats](
			n.config.URL+"/personal/rating/expected/json/",
			n.config.Timeout,
			nil,
		)
		errCtx := failure.Context{
			"url":         res.Request.URL,
			"status_code": strconv.Itoa(res.StatusCode),
			"body":        string(res.BodyByte),
		}
		if err != nil {
			return res, failure.Wrap(err, errCtx)
		}

		if res.StatusCode != http.StatusOK {
			return res, failure.New(apperr.NumbersAPIError, errCtx)
		}

		return res, nil
	}

	res, err := backoff.RetryWithData(operation, b)

	return res.Body.Data, failure.Wrap(err)
}
