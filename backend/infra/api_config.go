package infra

import "time"

type apiConfig struct {
	baseURL    string
	timeout    time.Duration
	retryCount int
}

func NewAPIConfig(
	baseURL string,
	timeout time.Duration,
	retryCount int,
) *apiConfig {
	return &apiConfig{
		baseURL:    baseURL,
		timeout:    timeout,
		retryCount: retryCount,
	}
}
