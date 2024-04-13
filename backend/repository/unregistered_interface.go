package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type UnregisteredInterface interface {
	Warship() (data.Warships, error)
}
