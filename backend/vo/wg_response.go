package vo

type WGError struct {
	Code    int
	Message string
	Field   string
	Value   string
}

type WGResponse interface {
	GetStatus() string
	GetError() WGError
}
