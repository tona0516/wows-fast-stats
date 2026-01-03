package infra

import "time"

type apiConfig struct {
	url        string
	retryCount int
	timeout    time.Duration
}

func NewApiConfig(
	url string,
	retryCount int,
	timeout time.Duration,
) *apiConfig {
	return &apiConfig{
		url:        url,
		retryCount: retryCount,
		timeout:    timeout,
	}
}
