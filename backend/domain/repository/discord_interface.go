package repository

type DiscordInterface interface {
	Comment(message string) error
}
