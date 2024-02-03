package webapi

type Request[T any] struct {
	URL    string
	Method string
	Body   T
}
