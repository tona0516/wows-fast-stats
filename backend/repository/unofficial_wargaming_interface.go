package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type UnofficialWargamingInterface interface {
	ClansAutoComplete(search string) (data.UWGClansAutocomplete, error)
}
