package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	//"github.com/samber/lo"
	"go-dianping/api"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
	"strconv"
	"time"
)

type ShopService interface {
	GetShopById(ctx context.Context, req *api.GetShopByIdReq) (*api.GetShopByIdRespData, error)
	UpdateShop(ctx context.Context, req *api.UpdateShopReq) error
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

func (s *shopService) GetShopById(ctx context.Context, params *api.GetShopByIdReq) (*api.GetShopByIdRespData, error) {
	var data api.GetShopByIdRespData

	// ========== check cache ==========
	key := constants.RedisCacheShopKey + params.Id
	cacheShopStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	var cacheShop *model.Shop
	if cacheShopStr != "" {
		if err := json.Unmarshal([]byte(cacheShopStr), &cacheShop); err != nil {
			return nil, err
		}
		if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
			return nil, err
		}
		return &data, nil
	}

	// ========== query sql db ==========
	id, err := strconv.Atoi(params.Id)
	if err != nil {
		return nil, err
	}
	shop, err := s.shopRepository.GetShopById(ctx, id)
	if err != nil {
		return nil, err
	}

	// ========== save to redis and return ==========
	ttl := time.Minute * constants.RedisCacheShopTTL
	bytes, err := json.Marshal(shop)
	if err != nil {
		return nil, err
	}
	if err := s.rdb.Set(ctx, key, bytes, ttl).Err(); err != nil {
		return nil, err
	}
	if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *shopService) UpdateShop(ctx context.Context, newShop *api.UpdateShopReq) error {
	// ========== update sql db ==========
	// already check id with gin.Context shouldBind in handler
	var shop model.Shop
	if err := copier.CopyWithOption(&shop, &newShop, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	if err := s.shopRepository.Update(ctx, &shop); err != nil {
		return err
	}

	// ========== delete cache ==========
	key := fmt.Sprintf("%s%d", constants.RedisLoginCodeKey, newShop.Id)
	s.rdb.Del(ctx, key)
	return nil
}
