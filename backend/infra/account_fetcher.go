package infra

import (
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/infra/webapi"

	"github.com/morikuni/failure"
)

type AccountFetcher struct {
	wargaming webapi.Wargaming
}

func NewAccountFetcher(
	wargaming webapi.Wargaming,
) *AccountFetcher {
	return &AccountFetcher{wargaming: wargaming}
}

func (f *AccountFetcher) Search(prefix string) (model.Accounts, error) {
	res, err := f.wargaming.AccountListForSearch(prefix)
	if err != nil {
		return nil, failure.Translate(err, apperr.FetchAccountListError)
	}

	result := make(model.Accounts)
	for _, v := range res {
		result[v.NickName] = v.AccountID
	}

	return result, nil
}

func (f *AccountFetcher) Fetch(playerNames []string) (model.Accounts, error) {
	res, err := f.wargaming.AccountList(playerNames)
	if err != nil {
		return nil, failure.Translate(err, apperr.FetchAccountListError)
	}

	result := make(model.Accounts)
	for _, v := range res {
		result[v.NickName] = v.AccountID
	}

	return result, nil
}
