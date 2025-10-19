package service

import (
	"context"
	"errors"
	"os"
	"wfs/backend/domain"
	"wfs/backend/repository"

	"github.com/morikuni/failure"
)

type BlackList struct {
	ctx          context.Context
	persistence  repository.PersistenceInterface
	notifyUpdate eventEmitFunc
}

func NewBlackList(
	ctx context.Context,
	persistence repository.PersistenceInterface,
	notifyUpdate eventEmitFunc,
) *BlackList {
	return &BlackList{
		ctx:          ctx,
		persistence:  persistence,
		notifyUpdate: notifyUpdate,
	}
}

func (b *BlackList) Get() (domain.BlackList, error) {
	list, err := b.persistence.LoadBlackList()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.BlackList{}, nil
		}

		return nil, failure.Wrap(err)
	}

	if list == nil {
		return domain.BlackList{}, nil
	}

	return list, nil
}

func (b *BlackList) Update(item domain.BlackListItem) error {
	list, err := b.Get()
	if err != nil {
		return failure.Wrap(err)
	}

	replaced := false
	for i := range list {
		if list[i].AccountID == item.AccountID {
			list[i] = item
			replaced = true
			break
		}
	}

	if !replaced {
		list = append(list, item)
	}

	if err := b.persistence.SaveBlackList(list); err != nil {
		return failure.Wrap(err)
	}

	b.notifyUpdate(b.ctx, EventUpdateBlackList, list)

	return nil
}

func (b *BlackList) Remove(accountID int) error {
	list, err := b.Get()
	if err != nil {
		return failure.Wrap(err)
	}

	filtered := make(domain.BlackList, 0, len(list))
	for _, item := range list {
		if item.AccountID != accountID {
			filtered = append(filtered, item)
		}
	}

	return failure.Wrap(b.persistence.SaveBlackList(filtered))
}
