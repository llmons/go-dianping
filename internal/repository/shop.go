package repository

import (
	"context"
	"go-dianping/internal/model"
)

type ShopRepository interface {
	GetShopById(ctx context.Context, id int) (*model.Shop, error)
	Update(ctx context.Context, shop *model.Shop) error
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

func (r *shopRepository) GetShopById(_ context.Context, id int) (*model.Shop, error) {
	var shop model.Shop
	r.db.First(&shop, id)
	return &shop, nil
}

func (r Repository) Update(_ context.Context, shop *model.Shop) error {
	r.db.Model(&model.Shop{}).Where("id = ?", shop.Model.Id).Updates(shop)
	return nil
}
