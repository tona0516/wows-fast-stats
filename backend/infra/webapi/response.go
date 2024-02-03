package webapi

type Response[T any, U any] struct {
	StatusCode int
	Request    Request[T]
	Body       U
	BodyByte   []byte
}
