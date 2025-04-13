package wargaming

type ResponseCommon[T any] struct {
	Status string `json:"status"`
	Error  Error  `json:"error"`
	Data   T      `json:"data"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

type Response interface {
	GetStatus() string
	GetError() Error
	Field() string
}

func (r ResponseCommon[T]) GetStatus() string {
	return r.Status
}

func (r ResponseCommon[T]) GetError() Error {
	return r.Error
}
