package infra

type APIClientInterface[T any] interface {
	GetRequest(query map[string]string) (T, error)
}
