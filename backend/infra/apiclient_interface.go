package infra

type APIClientInterface[T any] interface {
	GetRequest(query map[string]string) (APIResponse[T], error)
	PostMultipartFormData(forms []Form) (APIResponse[T], error)
}
