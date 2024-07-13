package infra

import (
	"net/http"
	"strconv"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
)

type Numbers struct {
	baseURL string
}

func NewNumbers(baseURL string) *Numbers {
	return &Numbers{baseURL: baseURL}
}

func (n *Numbers) ExpectedStats() (data.ExpectedStats, error) {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(2)

	var result data.NSExpectedStats
	resp, err := client.R().
		SetResult(&result).
		Get(n.baseURL + "/personal/rating/expected/json/")

	errCtx := failure.Context{
		"url":         resp.Request.URL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}
	if err != nil {
		return result.Data, failure.Wrap(err, errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return result.Data, failure.New(apperr.NumbersAPIError, errCtx)
	}

	return result.Data, failure.Wrap(err)
}
