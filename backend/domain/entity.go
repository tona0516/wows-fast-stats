package domain

type Entity[T primitive] interface {
	ID() T
	Equals(e *Entity[T]) bool
}
