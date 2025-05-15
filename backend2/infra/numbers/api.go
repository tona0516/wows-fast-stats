package numbers

import (
	"github.com/imroc/req/v3"
	"github.com/samber/do"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type API interface {
	Fetch() (Expected, error)
}

type api struct {
	client *req.Client
}

func NewAPI(i *do.Injector) (API, error) {
	return &api{client: do.MustInvokeNamed[*req.Client](i, "NumbersAPIClient")}, nil
}

func (a *api) Fetch() (Expected, error) {
	result := Expected{}

	_, err := a.client.R().
		SetSuccessResult(&result).
		Get("/personal/rating/expected/json/")

	return result, err
}
