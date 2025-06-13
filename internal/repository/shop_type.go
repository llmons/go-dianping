package repository

import (
	"go-dianping/internal/model"
)

type ShopTypeRepository interface {
	GetShopTypeList() ([]*model.ShopType, error)
}

func NewShopTypeRepository(
	repository *Repository,
) ShopTypeRepository {
	return &shopTypeRepository{
		Repository: repository,
	}
}

type shopTypeRepository struct {
	*Repository
}

func (r *shopTypeRepository) GetShopTypeList() ([]*model.ShopType, error) {
	var shopTypes []*model.ShopType
	err := r.db.Order("sort").Find(&shopTypes).Error
	if err != nil {
		return nil, err
	}
	return shopTypes, nil
}
