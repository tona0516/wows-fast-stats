package infra

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type APIClient[T any] struct{}

func (c *APIClient[T]) GetRequest(rawurl string) (T, error) {
	var response T

	u, err := url.Parse(rawurl)
	if err != nil {
		return response, err
	}

	res, err := http.Get(u.String())
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
