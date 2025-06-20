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
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
	"time"
)

type ShopTypeService interface {
	GetShopTypeList(ctx context.Context) (v1.GetShopTypeListRespData, error)
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

func (s *shopTypeService) GetShopTypeList(ctx context.Context) (v1.GetShopTypeListRespData, error) {
	// ========== check cache ==========
	cacheShopTypeStr, err := s.rdb.Get(ctx, constants.RedisCacheShopKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var cacheShopType v1.GetShopTypeListRespData
	if cacheShopTypeStr != "" {
		if err := json.Unmarshal([]byte(cacheShopTypeStr), &cacheShopType); err != nil {
			return nil, err
		}
		return cacheShopType, nil
	}

	// ========== query sql db ==========
	list, err := s.shopTypeRepository.GetShopTypeList(ctx)
	if err != nil {
		return v1.GetShopTypeListRespData{}, err
	}

	data := lo.Map(list, func(el *model.ShopType, idx int) *v1.GetShopTypeListRespDataItem {
		var item v1.GetShopTypeListRespDataItem
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
	ttl := time.Minute * constants.RedisCacheShopTypeTTL
	s.rdb.Set(ctx, constants.RedisCacheShopTypeKey, bytes, ttl)
	return data, nil
}
