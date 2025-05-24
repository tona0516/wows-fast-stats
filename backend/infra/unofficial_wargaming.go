package infra

import (
	"net/http"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type UnofficialWargaming struct {
	url        string
	maxRetry   int
	timeoutSec int
}

func NewUnofficialWargaming(
	url string,
	maxRetry int,
	timeoutSec int,
) *UnofficialWargaming {
	return &UnofficialWargaming{url: url, maxRetry: maxRetry, timeoutSec: timeoutSec}
}

func (w *UnofficialWargaming) ClansAutoComplete(search string) (data.UWGClansAutocomplete, error) {
	c := req.C().
		SetBaseURL(w.url).
		SetCommonRetryCount(w.maxRetry).
		SetTimeout(time.Duration(w.timeoutSec) * time.Second)

	var result data.UWGClansAutocomplete

	resp, err := c.R().
		SetSuccessResult(&result).
		SetQueryParams(map[string]string{
			"search": search,
			"type":   "clans",
		}).
		Get("/api/search/autocomplete/")
	if err != nil {
		return result, failure.Translate(err, apperr.UWGAPIError)
	}

	if resp.StatusCode != http.StatusOK {
		return result, failure.New(apperr.UWGAPIError, failure.Context{
			"status": resp.Status,
			"body":   resp.String(),
		})
	}

	return result, nil
}
