package repository

import (
	"context"
	"go-dianping/internal/entity"
)

type ShopTypeRepository interface {
	GetAll(ctx context.Context) ([]*entity.ShopType, error)
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

func (r *shopTypeRepository) GetAll(ctx context.Context) ([]*entity.ShopType, error) {
	return r.query.WithContext(ctx).ShopType.Order(r.query.ShopType.Sort.Asc()).Find()
}
