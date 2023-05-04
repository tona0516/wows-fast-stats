package infra

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type ApiClient[T any] struct {
}

func (c *ApiClient[T]) GetRequest(url string) (T, error) {
	var response T
	res, err := http.Get(url)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return response, errors.WithStack(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, errors.WithStack(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, errors.WithStack(err)
	}

	return response, nil
}
