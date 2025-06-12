package service

import (
	"context"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
)

type ShopTypeService interface {
	GetShopTypeList(ctx context.Context) ([]*model.ShopType, error)
}

func NewShopTypeService(
	service *Service,
	shopTypeRepository repository.ShopTypeRepository,
) ShopTypeService {
	return &shopTypeService{
		Service:            service,
		shopTypeRepository: shopTypeRepository,
	}
}

type shopTypeService struct {
	*Service
	shopTypeRepository repository.ShopTypeRepository
}

func (s *shopTypeService) GetShopTypeList(ctx context.Context) ([]*model.ShopType, error) {
	return s.shopTypeRepository.GetShopTypeList(ctx)
}
