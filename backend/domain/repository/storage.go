package repository

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type Storage interface {
	DataVersion() (uint, error)
	WriteDataVersion(version uint) error
}
