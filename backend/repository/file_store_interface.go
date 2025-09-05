package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type FileStoreInterface interface {
	Put(key data.FileStoreKey, value string) error
	Get(key data.FileStoreKey) (string, error)
	Delete(key data.FileStoreKey) error
	Keys() ([]string, error)
}
