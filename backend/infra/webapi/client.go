package webapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/morikuni/failure"
)

var ErrErrorResponse = errors.New("error response status code")

type ClientOption func(*Client)

type Client struct {
	baseURL    string
	path       string
	timeout    time.Duration
	headers    map[string]string
	query      map[string]string
	body       map[string]string
	isInsecure bool
}

func NewClient(baseURL string, options ...ClientOption) *Client {
	c := Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		timeout: 10 * time.Second,
	}

	for _, option := range options {
		option(&c)
	}

	return &c
}

func WithPath(path string) ClientOption {
	return func(rb *Client) {
		rb.path = strings.TrimPrefix(path, "/")
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(rb *Client) {
		rb.timeout = timeout
	}
}

func WithHeaders(headers map[string]string) ClientOption {
	return func(rb *Client) {
		rb.headers = headers
	}
}

func WithQuery(query map[string]string) ClientOption {
	return func(rb *Client) {
		rb.query = query
	}
}

func WithBody(body map[string]string) ClientOption {
	return func(rb *Client) {
		rb.body = body
	}
}

func WithIsInsecure(isInsecure bool) ClientOption {
	return func(rb *Client) {
		rb.isInsecure = isInsecure
	}
}

func (rb Client) GET() (http.Response, []byte, error) {
	return do(rb, makeGetRequest)
}

func (rb Client) POST() (http.Response, []byte, error) {
	return do(rb, makePostRequest)
}

func makeClient(rb Client) *http.Client {
	var client http.Client
	if rb.isInsecure {
		//nolint:gosec
		client = http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	} else {
		client = http.Client{}
	}

	client.Timeout = rb.timeout

	return &client
}

func makeGetRequest(rb Client) (*http.Request, error) {
	u, err := url.Parse(rb.baseURL)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	if len(rb.path) > 0 {
		u.Path = rb.path
	}

	q := u.Query()
	for k, v := range rb.query {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	for k, v := range rb.headers {
		req.Header.Add(k, v)
	}

	return req, nil
}

func makePostRequest(rb Client) (*http.Request, error) {
	b, err := json.Marshal(rb.body)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	u, err := url.Parse(rb.baseURL)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	if len(rb.path) > 0 {
		u.Path = rb.path
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, failure.Wrap(err)
	}

	for k, v := range rb.headers {
		req.Header.Add(k, v)
	}

	return req, nil
}

func do(rb Client, makeReqFunc func(rb Client) (*http.Request, error)) (http.Response, []byte, error) {
	var res http.Response
	var body []byte
	errCtx := make(failure.Context)

	client := makeClient(rb)
	req, err := makeReqFunc(rb)
	if err != nil {
		return res, body, failure.Wrap(err)
	}

	url := req.URL.Scheme + "://" + req.URL.Host
	if len(req.URL.Path) > 0 {
		url += req.URL.Path
	}
	if len(req.URL.RawQuery) > 0 {
		url += "?" + req.URL.RawQuery
	}

	errCtx["url"] = url

	_res, err := client.Do(req)
	if err != nil {
		return res, body, failure.Wrap(err, errCtx)
	}
	defer _res.Body.Close()
	res = *_res

	errCtx["status_code"] = strconv.Itoa(res.StatusCode)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return res, nil, failure.Wrap(err, errCtx)
	}

	errCtx["body"] = string(body)

	if isError(res) {
		return res, body, failure.Wrap(ErrErrorResponse, errCtx)
	}

	return res, body, failure.Wrap(err, errCtx)
}

func isError(res http.Response) bool {
	code := res.StatusCode
	return code > 399 && code < 600
}
