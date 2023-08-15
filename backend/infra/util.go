package infra

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"wfs/backend/application/vo"

	"github.com/morikuni/failure"
)

type APIResponse[T any] struct {
	FullURL    string
	StatusCode int
	Body       T
	ByteBody   []byte
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
	errCtx := failure.Context{}

	// build URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return response, failure.Wrap(err, errCtx)
	}
	q := u.Query()
	for _, query := range queries {
		q.Add(query.Key, query.Value)
	}
	u.RawQuery = q.Encode()
	response.FullURL = u.String()
	errCtx["url"] = u.String()

	// request
	res, err := http.Get(u.String())
	if err != nil {
		return response, failure.Wrap(err, errCtx)
	}
	defer res.Body.Close()
	response.StatusCode = res.StatusCode
	errCtx["status_code"] = strconv.Itoa(res.StatusCode)

	response.ByteBody, err = io.ReadAll(res.Body)
	if err != nil {
		return response, failure.Wrap(err, errCtx)
	}
	errCtx["body"] = string(response.ByteBody)

	// deserialize
	err = json.Unmarshal(response.ByteBody, &response.Body)
	if err != nil {
		return response, failure.Wrap(err, errCtx)
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
		return response, failure.Wrap(err)
	}

	// build request
	requestBody := &bytes.Buffer{}
	mw := multipart.NewWriter(requestBody)

	for _, form := range forms {
		//nolint:nestif
		if form.isFile {
			f, err := os.Open(form.content)
			if err != nil {
				return response, failure.Wrap(err)
			}

			fw, err := mw.CreateFormFile(form.name, form.content)
			if err != nil {
				return response, failure.Wrap(err)
			}

			_, err = io.Copy(fw, f)
			if err != nil {
				return response, failure.Wrap(err)
			}
		} else {
			if err := mw.WriteField(form.name, form.content); err != nil {
				return response, failure.Wrap(err)
			}
		}
	}

	// note: Close()で書き込まれる
	mw.Close()

	// request
	res, err := http.Post(u.String(), mw.FormDataContentType(), requestBody)
	if err != nil {
		return response, failure.Wrap(err)
	}
	defer res.Body.Close()
	response.StatusCode = res.StatusCode

	response.ByteBody, err = io.ReadAll(res.Body)
	if err != nil {
		return response, failure.Wrap(err)
	}

	// serialize
	err = json.Unmarshal(response.ByteBody, &response.Body)
	if err != nil {
		return response, failure.Wrap(err)
	}

	return response, nil
}
