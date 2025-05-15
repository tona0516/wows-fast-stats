package common

type Result[T any] struct {
	Value T
	Error error
}
