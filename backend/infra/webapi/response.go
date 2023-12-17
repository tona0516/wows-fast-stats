package webapi

type Response[T any] struct {
	FullURL    string
	StatusCode int
	Body       T
	ByteBody   []byte
}
