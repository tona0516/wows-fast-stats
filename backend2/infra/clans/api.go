package clans

import (
	"net/http"
	"wfs/backend/apperr"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/samber/do"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type API interface {
	FetchAutoComplete(search string) (Autocomplete, error)
}

type api struct {
	client *req.Client
}

func NewAPI(i *do.Injector) API {
	return &api{client: do.MustInvokeNamed[*req.Client](i, "ClansAPIClient")}
}

func (a *api) FetchAutoComplete(search string) (Autocomplete, error) {
	var result Autocomplete
	resp, err := a.client.R().
		SetQueryParams(map[string]string{
			"search": search,
			"type":   "clans",
		}).
		SetSuccessResult(&result).
		Get("/api/search/autocomplete/")

	if resp.StatusCode != http.StatusOK {
		return result, failure.New(apperr.UWGAPIError)
	}

	return result, failure.Translate(err, apperr.UWGAPIError)
}
