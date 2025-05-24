package infra

import (
	"encoding/json"
	"net/http"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type Numbers struct {
	url        string
	maxRetry   int
	timeoutSec int
}

func NewNumbers(
	url string,
	maxRetry int,
	timeoutSec int,
) *Numbers {
	return &Numbers{url: url, maxRetry: maxRetry, timeoutSec: timeoutSec}
}

func (n *Numbers) ExpectedStats() (data.ExpectedStats, error) {
	c := req.C().
		SetBaseURL(n.url).
		SetCommonRetryCount(n.maxRetry).
		SetTimeout(time.Duration(n.timeoutSec) * time.Second).
		EnableInsecureSkipVerify()

	resp, err := c.R().
		Get("/personal/rating/expected/json/")

	if err != nil {
		return nil, failure.Translate(err, apperr.NumbersAPIError)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, failure.New(apperr.NumbersAPIError, failure.Context{
			"status": resp.Status,
			"body":   resp.String(),
		})
	}

	var result data.NSExpectedStats
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, failure.Translate(err, apperr.NumbersAPIError)
	}

	return result.Data, nil
}
