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

type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) Debug(msg string) {
	l.Called(msg)
}

func (l *mockLogger) Info(msg string) {
	l.Called(msg)
}

func (l *mockLogger) Warn(msg string, err error) {
	l.Called(msg, err)
}

func (l *mockLogger) Error(msg string, err error) {
	l.Called(msg, err)
}
