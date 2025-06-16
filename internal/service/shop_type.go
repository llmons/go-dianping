package service

import (
	"go-dianping/api"
	"go-dianping/internal/repository"
)

type ShopTypeService interface {
	GetShopTypeList() (api.GetShopTypeListRespData, error)
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

func (s *shopTypeService) GetShopTypeList() (api.GetShopTypeListRespData, error) {
	list, err := s.shopTypeRepository.GetShopTypeList()
	if err != nil {
		return api.GetShopTypeListRespData{}, err
	}

	data := make(api.GetShopTypeListRespData, len(list))
	for i, shopType := range list {
		data[i] = &api.GetShopTypeListRespDataItem{
			Id:   shopType.Id,
			Name: shopType.Name,
			Icon: shopType.Icon,
			Sort: shopType.Sort,
		}
	}
	return data, nil
}
