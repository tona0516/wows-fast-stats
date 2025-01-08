package infra

import (
	"encoding/json"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra/webapi"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
)

type UnofficialWargaming struct {
	config RequestConfig
}

func NewUnofficialWargaming(config RequestConfig) *UnofficialWargaming {
	return &UnofficialWargaming{config: config}
}

func (w *UnofficialWargaming) ClansAutoComplete(search string) (data.UWGClansAutocomplete, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), w.config.Retry)
	operation := func() (data.UWGClansAutocomplete, error) {
		var result data.UWGClansAutocomplete

		_, body, err := webapi.NewClient(w.config.URL,
			webapi.WithPath("/api/search/autocomplete/"),
			webapi.WithQuery(map[string]string{
				"search": search,
				"type":   "clans",
			}),
			webapi.WithTimeout(w.config.Timeout),
			webapi.WithIsInsecure(true),
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
		return res, failure.Translate(err, apperr.UWGAPIError)
	}

	return res, nil
}
