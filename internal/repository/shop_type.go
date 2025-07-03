package repository

import (
	"context"
	"go-dianping/internal/model"
)

type ShopTypeRepository interface {
	GetAll(ctx context.Context) ([]*model.ShopType, error)
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

func (r *shopTypeRepository) GetAll(ctx context.Context) ([]*model.ShopType, error) {
	return r.query.WithContext(ctx).ShopType.Order(r.query.ShopType.Sort.Asc()).Find()
}
