package repository

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type DiscordInterface interface {
	Comment(message string) error
}
