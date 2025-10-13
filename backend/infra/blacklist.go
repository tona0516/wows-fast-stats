package infra

import (
	"wfs/backend/domain"
)

type BlackList struct {
	basePath string
}

func NewBlackList(basePath string) *BlackList {
	return &BlackList{
		basePath: basePath,
	}
}

func (b *BlackList) Load() (*domain.BlackList, error) {
	return load[domain.BlackList](b.basePath)
}

func (b *BlackList) Save(data domain.BlackList) error {
	return save(b.basePath, data)
}
