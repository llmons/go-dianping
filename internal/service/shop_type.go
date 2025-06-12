package service

import (
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
)

type ShopTypeService interface {
	GetShopTypeList() ([]*model.ShopType, error)
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

func (s *shopTypeService) GetShopTypeList() ([]*model.ShopType, error) {
	return s.shopTypeRepository.GetShopTypeList()
}
