package domain

import (
	"encoding/json"
	"fmt"
)

type ValueObject2[T primitive, U primitive] struct {
	first  T
	second U
}

func (v ValueObject2[T, U]) Value() (T, U) {
	return v.first, v.second
}

func (v ValueObject2[T, U]) String() string {
	return fmt.Sprintf("%v %v", v.first, v.second)
}

func (v ValueObject2[T, U]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}
