package repository

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type DiscordInterface interface {
	Comment(message string) error
}
