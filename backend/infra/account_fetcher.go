package infra

import (
	"encoding/json"
	"strings"
	"wfs/backend/domain/model"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type AccountFetcher struct {
	wargamingClient req.Client
}

func NewAccountFetcher(
	wargamingClient req.Client,
) *AccountFetcher {
	return &AccountFetcher{wargamingClient: wargamingClient}
}

func (f *AccountFetcher) Search(prefix string) (model.Accounts, error) {
	var rb WGAccountListResponse
	resp, err := f.wargamingClient.R().
		AddQueryParam("search", prefix).
		AddQueryParam("fields", WGAccountListResponse{}.Field()).
		AddQueryParam("limit", "10").
		Get("/wows/account/list/")
	if err != nil {
		return nil, failure.Wrap(err)
	}

	if err := json.Unmarshal(resp.Bytes(), &rb); err != nil {
		return nil, failure.Wrap(err)
	}

	result := make(model.Accounts)
	for _, v := range rb.Data {
		result[v.NickName] = v.AccountID
	}

	return result, nil
}

func (f *AccountFetcher) Fetch(playerNames []string) (model.Accounts, error) {
	var rb WGAccountListResponse
	resp, err := f.wargamingClient.R().
		AddQueryParam("search", strings.Join(playerNames, ",")).
		AddQueryParam("fields", WGAccountListResponse{}.Field()).
		AddQueryParam("type", "exact").
		Get("/wows/account/list/")
	if err != nil {
		return nil, failure.Wrap(err)
	}

	if err := json.Unmarshal(resp.Bytes(), &rb); err != nil {
		return nil, failure.Wrap(err)
	}

	result := make(model.Accounts)
	for _, v := range rb.Data {
		result[v.NickName] = v.AccountID
	}

	return result, nil
}
