package repository

//go:generate mockgen -source=$GOFILE -destination ../mock_$GOPACKAGE/$GOFILE -package mock_$GOPACKAGE
type DiscordInterface interface {
	Comment(message string) error
}
