package infra

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/cenkalti/backoff/v4"
)

type APIClient[T any] struct{}

func (c *APIClient[T]) GetRequest(rawurl string) (T, error) {
	var response T

	u, err := url.Parse(rawurl)
	if err != nil {
		return response, err
	}

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
	operation := func() (*http.Response, error) {
		return http.Get(u.String())
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return response, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
