package domain

type WGError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

type WGResponse interface {
	GetStatus() string
	GetError() WGError
}
