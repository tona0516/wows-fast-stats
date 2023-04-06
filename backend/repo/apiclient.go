package repo

import (
	"encoding/json"
	"io"
	"net/http"
)

type ApiClient[T any] struct {
}

func (a *ApiClient[T]) GetRequest(url string) (T, error) {
	var response T
	res, err := http.Get(url)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return response, err
	}

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
