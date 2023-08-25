package response

type WGResponseCommon[T any] struct {
	Status string  `json:"status"`
	Error  WGError `json:"error"`
	Data   T       `json:"data"`
}

type WGError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

type WGResponse interface {
	GetStatus() string
	GetError() WGError
	Field() string
}

func (r WGResponseCommon[T]) GetStatus() string {
	return r.Status
}

func (r WGResponseCommon[T]) GetError() WGError {
	return r.Error
}
