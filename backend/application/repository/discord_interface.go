package repository

type DiscordInterface interface {
	Upload(text string, message string) error
}
