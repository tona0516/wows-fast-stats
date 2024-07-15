package response

type WGResponse[T any] struct {
	Status string `json:"status"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Data T `json:"data"`
}
