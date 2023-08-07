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
	"wfs/backend/application/vo"
)

type APIResponse[T any] struct {
	StatusCode int
	BodyByte   []byte
	Body       T
}

type Form struct {
	name    string
	content string
	isFile  bool
}

func getRequest[T any](
	rawURL string,
	queries ...vo.Pair,
) (APIResponse[T], error) {
	var response APIResponse[T]

	// build URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
	}
	q := u.Query()
	for _, query := range queries {
		q.Add(query.Key, query.Value)
	}
	u.RawQuery = q.Encode()

	// request
	res, err := http.Get(u.String())
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
	}

	response.StatusCode = res.StatusCode

	defer res.Body.Close()

	response.BodyByte, err = io.ReadAll(res.Body)
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
	}

	// serialize
	err = json.Unmarshal(response.BodyByte, &response.Body)
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
	}

	return response, nil
}

//nolint:cyclop
func postMultipartFormData[T any](
	rawURL string,
	forms []Form,
) (APIResponse[T], error) {
	var response APIResponse[T]

	// build URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
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
	res, err := http.Post(u.String(), mw.FormDataContentType(), requestBody)
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
	}

	response.StatusCode = res.StatusCode

	defer res.Body.Close()

	response.BodyByte, err = io.ReadAll(res.Body)
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
	}

	// serialize
	err = json.Unmarshal(response.BodyByte, &response.Body)
	if err != nil {
		return response, apperr.New(apperr.ErrHTTPRequest, err)
	}

	return response, nil
}
