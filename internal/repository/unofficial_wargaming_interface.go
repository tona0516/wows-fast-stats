package repository

import "wfs/internal/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type UnofficialWargamingInterface interface {
	ClansAutoComplete(search string) (data.UWGClansAutocomplete, error)
}
