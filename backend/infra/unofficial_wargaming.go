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

type UnofficialWargaming struct {
	baseURL string
}

func NewUnofficialWargaming(baseURL string) *UnofficialWargaming {
	return &UnofficialWargaming{baseURL: baseURL}
}

func (w *UnofficialWargaming) ClansAutoComplete(search string) (data.UWGClansAutocomplete, error) {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(2)

	var result data.UWGClansAutocomplete
	resp, err := client.R().
		SetResult(&result).
		SetQueryParams(map[string]string{
			"search": search,
			"type":   "clans",
		}).
		Get(w.baseURL + "/api/search/autocomplete/")

	errCtx := failure.Context{
		"url":         resp.Request.URL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}
	if err != nil {
		return result, failure.Wrap(err, errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return result, failure.New(apperr.UWGAPIError, errCtx)
	}

	return result, failure.Wrap(err)
}
