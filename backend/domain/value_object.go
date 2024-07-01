package domain

import (
	"encoding/json"
	"fmt"
)

type ValueObject[T primitive] struct {
	value T
}

func (v ValueObject[T]) Value() T {
	return v.value
}

func (v ValueObject[T]) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v ValueObject[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}
