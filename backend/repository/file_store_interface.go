package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type FileStoreInterface interface {
	Put(key data.FileStorePath, value string) error
	Get(key data.FileStorePath) (string, error)
	Delete(key data.FileStorePath) error
	Files(childPath string) ([]string, error)
}
