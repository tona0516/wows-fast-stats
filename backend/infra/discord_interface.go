package infra

type DiscordInterface interface {
	Upload(text string) (APIResponse[any], error)
}
