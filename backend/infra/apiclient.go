package infra

import (
	"changeme/backend/apperr"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"unsafe"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

type APIClient[T any] struct {
	baseURL string
	logger  LoggerInterface
}

func NewAPIClient[T any](
	baseURL string,
	logger LoggerInterface,
) *APIClient[T] {
	return &APIClient[T]{
		baseURL: baseURL,
		logger:  logger,
	}
}

func (c *APIClient[T]) GetRequest(query map[string]string) (T, error) {
	var response T
	var statusCode int
	var bodyString string

	// build URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		c.logError(c.baseURL, statusCode, bodyString, err)
		return response, errors.WithStack(apperr.Ac.Parse.WithRaw(err))
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
		c.logError(u.String(), statusCode, bodyString, err)
		return response, errors.WithStack(apperr.Ac.Retry.WithRaw(err))
	}

	statusCode = res.StatusCode

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logError(u.String(), statusCode, bodyString, err)
		return response, errors.WithStack(apperr.Ac.Read.WithRaw(err))
	}

	bodyString = *(*string)(unsafe.Pointer(&body))

	// serialize
	err = json.Unmarshal(body, &response)
	if err != nil {
		c.logError(u.String(), statusCode, bodyString, err)
		return response, errors.WithStack(apperr.Ac.Unmarshal.WithRaw(err))
	}

	return response, nil
}

func (c *APIClient[T]) logError(url string, statusCode int, responseBody string, err error) {
	c.logger.Error(fmt.Sprintf("URL: %s, STATUS_CODE: %d RESPONSE: %s", url, statusCode, responseBody), err)
}
