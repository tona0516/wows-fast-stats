package webapi

import (
	"encoding/json"
	"wfs/backend/apperr"
	"wfs/backend/infra/response"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type UnofficialWargaming interface {
	ClansAutoComplete(search string) (response.UWGClansAutocomplete, error)
}

type unofficialWargaming struct {
	config RequestConfig
}

func NewUnofficialWargaming(config RequestConfig) UnofficialWargaming {
	return &unofficialWargaming{config: config}
}

func (w *unofficialWargaming) ClansAutoComplete(search string) (response.UWGClansAutocomplete, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), w.config.Retry)
	operation := func() (response.UWGClansAutocomplete, error) {
		var result response.UWGClansAutocomplete

		_, body, err := NewClient(w.config.URL,
			WithPath("/api/search/autocomplete/"),
			WithQuery(map[string]string{
				"search": search,
				"type":   "clans",
			}),
			WithTimeout(w.config.Timeout),
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
