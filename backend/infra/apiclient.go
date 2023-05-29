package infra

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/cenkalti/backoff/v4"
)

type APIClient[T any] struct {
	baseURL string
}

func NewAPIClient[T any](baseURL string) *APIClient[T] {
	return &APIClient[T]{baseURL: baseURL}
}

func (c *APIClient[T]) GetRequest(query map[string]string) (T, error) {
	var response T

	// build URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return response, err
	}
	q := u.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	// request
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
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

	// serialize
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
