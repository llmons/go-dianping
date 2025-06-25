package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/entity"
	"go-dianping/internal/repository"
)

type ShopTypeService interface {
	QueryTypeList(ctx context.Context) (v1.QueryTypeListRespData, error)
}

type shopTypeService struct {
	*Service
	shopTypeRepository repository.ShopTypeRepository
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

func (s *shopTypeService) QueryTypeList(ctx context.Context) (v1.QueryTypeListRespData, error) {
	// ========== check cache ==========
	cacheShopTypeStr, err := s.rdb.Get(ctx, constants.RedisCacheShopKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var cacheShopType v1.QueryTypeListRespData
	if cacheShopTypeStr != "" {
		if err := json.Unmarshal([]byte(cacheShopTypeStr), &cacheShopType); err != nil {
			return nil, err
		}
		return cacheShopType, nil
	}

	// ========== query sql db ==========
	list, err := s.shopTypeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	data := lo.Map(list, func(el *entity.ShopType, idx int) *v1.QueryTypeListRespDataItem {
		var item v1.QueryTypeListRespDataItem
		if err := copier.Copy(&item, el); err != nil {
			return nil
		}
		return &item
	})

	// ========== save to redis and return ==========
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	s.rdb.Set(ctx, constants.RedisCacheShopTypeKey, bytes, constants.RedisCacheShopTypeTTL)
	return data, nil
}
