package infra

import (
	"github.com/stretchr/testify/mock"
)

type mockAPIClient[T any] struct {
	mock.Mock
}

func (m *mockAPIClient[T]) GetRequest(_ map[string]string) (APIResponse[T], error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(APIResponse[T]), args.Error(1)
}

func (m *mockAPIClient[T]) PostMultipartFormData(_ []Form) (APIResponse[T], error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(APIResponse[T]), args.Error(1)
}
