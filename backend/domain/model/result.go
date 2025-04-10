package model

type Result[T any] struct {
	Value T
	Error error
}
