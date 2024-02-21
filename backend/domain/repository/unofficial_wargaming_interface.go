package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock_$GOPACKAGE/$GOFILE -package mock_$GOPACKAGE
type UnofficialWargamingInterface interface {
	ClansAutoComplete(search string) (model.UWGClansAutocomplete, error)
}
