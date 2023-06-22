package infra

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"wfs/backend/apperr"

	"unsafe"

	"github.com/cenkalti/backoff/v4"
)

type APIClient[T any] struct {
	baseURL string
}

type APIResponse[T any] struct {
	StatusCode int
	Body       T
	BodyString string
}

type Form struct {
	name    string
	content string
	isFile  bool
}

func NewAPIClient[T any](
	baseURL string,
) *APIClient[T] {
	return &APIClient[T]{baseURL: baseURL}
}

func (c *APIClient[T]) GetRequest(query map[string]string) (APIResponse[T], error) {
	var response APIResponse[T]

	// build URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
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
		return response, apperr.New(apperr.HTTPRequest, err)
	}

	response.StatusCode = res.StatusCode

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
	}

	response.BodyString = *(*string)(unsafe.Pointer(&body))

	// serialize
	err = json.Unmarshal(body, &response.Body)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
	}

	return response, nil
}

//nolint:cyclop
func (c *APIClient[T]) PostMultipartFormData(forms []Form) (APIResponse[T], error) {
	var response APIResponse[T]

	// build URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
	}

	// build request
	requestBody := &bytes.Buffer{}
	mw := multipart.NewWriter(requestBody)

	for _, form := range forms {
		//nolint:nestif
		if form.isFile {
			f, err := os.Open(form.content)
			if err != nil {
				return response, err
			}

			fw, err := mw.CreateFormFile(form.name, form.content)
			if err != nil {
				return response, err
			}

			_, err = io.Copy(fw, f)
			if err != nil {
				return response, err
			}
		} else {
			if err := mw.WriteField(form.name, form.content); err != nil {
				return response, err
			}
		}
	}

	// note: Close()で書き込まれる
	mw.Close()

	// request
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	operation := func() (*http.Response, error) {
		return http.Post(u.String(), mw.FormDataContentType(), requestBody)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
	}

	response.StatusCode = res.StatusCode

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
	}

	response.BodyString = *(*string)(unsafe.Pointer(&body))

	// serialize
	err = json.Unmarshal(body, &response.Body)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
	}

	return response, nil
}
