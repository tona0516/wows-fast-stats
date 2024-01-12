package repository

import "wfs/backend/domain/model"

type UnofficialWargamingInterface interface {
	ClansAutoComplete(search string) (model.UWGClansAutocomplete, error)
}
