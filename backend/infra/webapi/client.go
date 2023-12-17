package webapi

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"wfs/backend/application/vo"

	"github.com/morikuni/failure"
)

func GetRequest[T any](
	rawURL string,
	timeout time.Duration,
	queries ...vo.Pair,
) (Response[T], error) {
	var response Response[T]
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
	client := http.Client{}
	client.Timeout = timeout

	res, err := client.Get(u.String())
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
func PostMultipartFormData[T any](
	rawURL string,
	timeout time.Duration,
	forms []Form,
) (Response[T], error) {
	var response Response[T]
	errCtx := failure.Context{}

	// build URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return response, failure.Wrap(err, errCtx)
	}
	response.FullURL = u.String()
	errCtx["url"] = u.String()

	// build request
	requestBody := &bytes.Buffer{}
	mw := multipart.NewWriter(requestBody)

	for _, form := range forms {
		//nolint:nestif
		if form.isFile {
			f, err := os.Open(form.content)
			if err != nil {
				return response, failure.Wrap(err, errCtx)
			}

			fw, err := mw.CreateFormFile(form.name, form.content)
			if err != nil {
				return response, failure.Wrap(err, errCtx)
			}

			_, err = io.Copy(fw, f)
			if err != nil {
				return response, failure.Wrap(err, errCtx)
			}
		} else {
			if err := mw.WriteField(form.name, form.content); err != nil {
				return response, failure.Wrap(err, errCtx)
			}
		}
	}

	// note: Close()で書き込まれる
	mw.Close()

	// request
	client := http.Client{}
	client.Timeout = timeout

	res, err := client.Post(u.String(), mw.FormDataContentType(), requestBody)
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
