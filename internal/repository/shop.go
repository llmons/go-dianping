package repository

import (
	"context"
	"go-dianping/internal/entity"
	"gorm.io/gen"
)

type ShopRepository interface {
	GetById(ctx context.Context, id int64) (*entity.Shop, error)
	Updates(ctx context.Context, shop *entity.Shop) (gen.ResultInfo, error)
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

func (r *shopRepository) GetById(ctx context.Context, id int64) (*entity.Shop, error) {
	return r.query.WithContext(ctx).Shop.Where(r.query.Shop.ID.Eq(id)).First()
}

func (r Repository) Updates(ctx context.Context, shop *entity.Shop) (gen.ResultInfo, error) {
	return r.query.WithContext(ctx).Shop.Where(r.query.Shop.ID.Eq(shop.ID)).Updates(shop)
}
