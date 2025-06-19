package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-dianping/api"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/repository"
	"time"
)

type ShopTypeService interface {
	GetShopTypeList(ctx context.Context) (api.GetShopTypeListRespData, error)
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

func (s *shopTypeService) GetShopTypeList(ctx context.Context) (api.GetShopTypeListRespData, error) {
	// ========== check cache ==========
	cacheShopTypeStr, err := s.rdb.Get(ctx, constants.RedisCacheShopKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	fmt.Println("DEBUG")

	var cacheShopType api.GetShopTypeListRespData
	if cacheShopTypeStr != "" {
		if err := json.Unmarshal([]byte(cacheShopTypeStr), &cacheShopType); err != nil {
			return nil, err
		}
		return cacheShopType, nil
	}

	// ========== query sql db ==========
	list, err := s.shopTypeRepository.GetShopTypeList(ctx)
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

	// ========== save to redis and return ==========
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	ttl := time.Minute * constants.RedisCacheShopTypeTTL
	s.rdb.Set(ctx, constants.RedisCacheShopTypeKey, bytes, ttl)
	return data, nil
}
