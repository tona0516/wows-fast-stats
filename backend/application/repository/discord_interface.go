package repository

type DiscordInterface interface {
	Upload(text string) error
}
