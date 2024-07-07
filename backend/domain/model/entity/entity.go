package entity

type Entity[T comparable] interface {
	ID() T
	Equals(e *Entity[T]) bool
}
