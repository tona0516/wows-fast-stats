package vo

type Result[T any] struct {
	Value T
	Error error
}
