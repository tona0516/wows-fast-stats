package infra

import (
	"net/http"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
)

type UnofficialWargaming struct {
	config apiConfig
}

func NewUnofficialWargaming(config apiConfig) *UnofficialWargaming {
	return &UnofficialWargaming{config: config}
}

func (w *UnofficialWargaming) ClansAutoComplete(search string) (data.UWGClansAutocomplete, error) {
	client := resty.New().
		SetTimeout(w.config.timeout).
		SetRetryCount(w.config.retryCount)

	var result data.UWGClansAutocomplete
	resp, err := client.R().
		SetResult(&result).
		SetQueryParams(map[string]string{
			"search": search,
			"type":   "clans",
		}).
		Get(w.config.baseURL + "/api/search/autocomplete/")

	errCtx := failure.Context{
		"url":         resp.Request.URL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}

	if err != nil {
		return result, failure.New(apperr.UWGAPIError, failure.Messagef(err.Error()), errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return result, failure.New(apperr.UWGAPIError, errCtx)
	}

	return result, nil
}
