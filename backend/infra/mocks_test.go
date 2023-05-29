package infra

import (
	"github.com/stretchr/testify/mock"
)

type mockAPIClient[T any] struct {
	mock.Mock
}

func (m *mockAPIClient[T]) GetRequest(_ map[string]string) (T, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(T), args.Error(1)
}
