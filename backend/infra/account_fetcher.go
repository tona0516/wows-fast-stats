package infra

import (
	"strings"
	"wfs/backend/domain/model"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type AccountFetcher struct {
	wargamingClient *req.Client
}

func NewAccountFetcher(wargamingClient *req.Client) *AccountFetcher {
	return &AccountFetcher{wargamingClient: wargamingClient}
}

func (f *AccountFetcher) FetchByPrefix(prefix string) (model.Accounts, error) {
	accounts, err := f.get(map[string]string{
		"search": prefix,
		"limit":  "10",
	})
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return accounts, nil
}

func (f *AccountFetcher) FetchByNames(names []string) (model.Accounts, error) {
	accounts, err := f.get(map[string]string{
		"search": strings.Join(names, ","),
		"type":   "exact",
	})
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return accounts, nil
}

func (f *AccountFetcher) get(additionalParams map[string]string) (model.Accounts, error) {
	var body WGAccountListResponse

	req := f.wargamingClient.R().
		SetSuccessResult(&body).
		AddQueryParam("fields", WGAccountListResponse{}.Field())

	for k, v := range additionalParams {
		req.AddQueryParam(k, v)
	}

	_, err := req.Get("/wows/account/list/")
	if err != nil {
		return nil, failure.Wrap(err)
	}

	result := make(model.Accounts)
	for _, v := range body.Data {
		result[v.NickName] = v.AccountID
	}

	return result, nil
}
