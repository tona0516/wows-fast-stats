package webapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/morikuni/failure"
)

func GetRequest[T any](
	rawURL string,
	timeout time.Duration,
	queries map[string]string,
) (Response[any, T], error) {
	var response Response[any, T]
	response.Request.Method = http.MethodGet

	// build URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return response, failure.Wrap(err)
	}
	q := u.Query()
	for key, value := range queries {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()
	response.Request.URL = u.String()

	// build request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return response, failure.Wrap(err)
	}

	// request
	err = request[any, T](req, &response, timeout)
	return response, failure.Wrap(err)
}

func PostRequestJSON[T, U any](
	rawURL string,
	timeout time.Duration,
	requestBody T,
) (Response[T, U], error) {
	var response Response[T, U]
	response.Request.Method = http.MethodPost
	response.Request.URL = rawURL
	response.Request.Body = requestBody

	// unmarshal body to json
	requestBodyByte, err := json.Marshal(requestBody)
	if err != nil {
		return response, failure.Wrap(err)
	}

	// build request
	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewBuffer(requestBodyByte))
	if err != nil {
		return response, failure.Wrap(err)
	}
	req.Header.Set("content-type", "application/json")

	// request
	err = request[T, U](req, &response, timeout)
	return response, failure.Wrap(err)
}

func request[T, U any](
	request *http.Request,
	response *Response[T, U],
	timeout time.Duration,
) error {
	client := http.Client{}
	client.Timeout = timeout

	// request
	res, err := client.Do(request)
	if err != nil {
		return failure.Wrap(err)
	}
	defer res.Body.Close()

	// read response
	response.StatusCode = res.StatusCode
	response.BodyByte, err = io.ReadAll(res.Body)
	if err != nil {
		return failure.Wrap(err)
	}

	// deserialize
	err = json.Unmarshal(response.BodyByte, &response.Body)
	if err != nil {
		return failure.Wrap(err)
	}

	return nil
}
