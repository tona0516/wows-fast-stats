package infra

import (
	"net/http"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
)

type Numbers struct {
	config apiConfig
}

func NewNumbers(config apiConfig) *Numbers {
	return &Numbers{config: config}
}

func (n *Numbers) ExpectedStats() (data.ExpectedStats, error) {
	client := resty.New().
		SetTimeout(n.config.timeout).
		SetRetryCount(n.config.retryCount)

	var result data.NSExpectedStats
	resp, err := client.R().
		SetResult(&result).
		Get(n.config.baseURL + "/personal/rating/expected/json/")

	errCtx := failure.Context{
		"url":         resp.Request.URL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}

	if err != nil {
		return result.Data, failure.New(apperr.NumbersAPIError, failure.Messagef(err.Error()), errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return result.Data, failure.New(apperr.NumbersAPIError, errCtx)
	}

	return result.Data, nil
}
