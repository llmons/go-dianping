package service

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/entity"
	"go-dianping/internal/repository"
)

type ShopTypeService interface {
	QueryTypeList(ctx context.Context) ([]v1.QueryTypeListRespDataItem, error)
}

type shopTypeService struct {
	*Service
	shopTypeRepo repository.ShopTypeRepository
}

func NewShopTypeService(
	service *Service,
	shopTypeRepository repository.ShopTypeRepository,
) ShopTypeService {
	return &shopTypeService{
		Service:      service,
		shopTypeRepo: shopTypeRepository,
	}
}

func (s *shopTypeService) QueryTypeList(ctx context.Context) ([]v1.QueryTypeListRespDataItem, error) {
	// 使用 List 结构

	exist, err := s.rdb.Exists(ctx, constants.RedisCacheShopTypeKey).Result()
	if err != nil {
		return nil, err
	}
	if exist == 1 {
		strSlices, err := s.rdb.LRange(ctx, constants.RedisCacheShopTypeKey, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		return lo.Map(strSlices, func(el string, idx int) v1.QueryTypeListRespDataItem {
			var item v1.QueryTypeListRespDataItem
			if err := json.Unmarshal([]byte(el), &item); err != nil {
				return item
			}
			return item
		}), nil
	}

	list, err := s.shopTypeRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	data := lo.Map(list, func(el *entity.ShopType, idx int) v1.QueryTypeListRespDataItem {
		var item v1.QueryTypeListRespDataItem
		if err := copier.Copy(&item, el); err != nil {
			return item
		}

		bytes, err := json.Marshal(item)
		if err != nil {
			return v1.QueryTypeListRespDataItem{}
		}
		if err := s.rdb.RPush(ctx, constants.RedisCacheShopTypeKey, string(bytes)).Err(); err != nil {
			return v1.QueryTypeListRespDataItem{}
		}

		return item
	})

	return data, nil
}
