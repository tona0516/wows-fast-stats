package infra

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"unsafe"
	"wfs/backend/apperr"

	"github.com/cenkalti/backoff/v4"
)

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

func getRequest[T any](
	rawURL string,
	query map[string]string,
	retry uint64,
) (APIResponse[T], error) {
	var response APIResponse[T]

	// build URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return response, apperr.New(apperr.HTTPRequest, err)
	}
	q := u.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	// request
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retry)
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
func postMultipartFormData[T any](
	rawURL string,
	forms []Form,
	retry uint64,
) (APIResponse[T], error) {
	var response APIResponse[T]

	// build URL
	u, err := url.Parse(rawURL)
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
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retry)
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

func readJSON[T any](path string, defaultValue T) (T, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return defaultValue, err
	}

	if err := json.Unmarshal(f, &defaultValue); err != nil {
		return defaultValue, err
	}

	return defaultValue, nil
}

func writeJSON[T any](path string, target T) error {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	return encoder.Encode(target)
}
