package repository

import (
	"context"
	"go-dianping/internal/model"
	"gorm.io/gen"
)

type ShopRepository interface {
	GetById(ctx context.Context, id uint64) (*model.Shop, error)
	Updates(ctx context.Context, shop *model.Shop) (gen.ResultInfo, error)
}

func NewShopRepository(
	repository *Repository,
) ShopRepository {
	return &shopRepository{
		Repository: repository,
	}
}

type shopRepository struct {
	*Repository
}

func (r *shopRepository) GetById(ctx context.Context, id uint64) (*model.Shop, error) {
	return r.query.WithContext(ctx).Shop.Where(r.query.Shop.ID.Eq(id)).First()
}

func (r *shopRepository) Updates(ctx context.Context, shop *model.Shop) (gen.ResultInfo, error) {
	return r.query.WithContext(ctx).Shop.Where(r.query.Shop.ID.Eq(shop.ID)).Updates(shop)
}
