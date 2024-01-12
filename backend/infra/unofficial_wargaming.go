package infra

import (
	"net/http"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
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

func (w *UnofficialWargaming) ClansAutoComplete(search string) (model.UWGClansAutocomplete, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), w.config.Retry)
	operation := func() (webapi.Response[model.UWGClansAutocomplete], error) {
		res, err := webapi.GetRequest[model.UWGClansAutocomplete](
			w.config.URL+"/api/search/autocomplete/",
			w.config.Timeout,
			map[string]string{
				"search": search,
				"type":   "clans",
			},
		)
		if err != nil {
			return res, err
		}

		if res.StatusCode != http.StatusOK {
			return res, failure.New(apperr.UWGAPIError)
		}

		return res, nil
	}

	res, err := backoff.RetryWithData(operation, b)

	return res.Body, failure.Wrap(err, apperr.ToRequestErrorContext(res))
}
