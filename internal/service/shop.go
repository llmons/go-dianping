package service

import (
	"context"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
)

type ShopService interface {
	GetShopById(ctx context.Context, id int64) (*model.Shop, error)
}

func NewShopService(
	service *Service,
	shopRepository repository.ShopRepository,
) ShopService {
	return &shopService{
		Service:        service,
		shopRepository: shopRepository,
	}
}

type shopService struct {
	*Service
	shopRepository repository.ShopRepository
}

func (s *shopService) GetShopById(ctx context.Context, id int64) (*model.Shop, error) {
	return s.shopRepository.GetShopById(ctx, id)
}
