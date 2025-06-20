package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"
	"gorm.io/gorm"

	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
	"strconv"
	"time"
)

type ShopService interface {
	GetShopById(ctx context.Context, req *v1.GetShopByIdReq) (*v1.GetShopByIdRespData, error)
	UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error
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

func (s *shopService) GetShopById(ctx context.Context, req *v1.GetShopByIdReq) (*v1.GetShopByIdRespData, error) {
	var data v1.GetShopByIdRespData

	// ========== check cache ==========
	key := constants.RedisCacheShopKey + req.Id
	cacheShopStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// err == nil || err == redis.Nil
	if err == nil && cacheShopStr == "" {
		return nil, errors.New("shop not exist")
	}
	if cacheShopStr != "" {
		var cacheShop model.Shop
		if err := json.Unmarshal([]byte(cacheShopStr), &cacheShop); err != nil {
			return nil, err
		}
		if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
			return nil, err
		}
		return &data, nil
	}

	// ========== query sql db ==========
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}
	shop, err := s.shopRepository.GetShopById(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// solve the cache penetration
			s.rdb.Set(ctx, key, "", time.Minute*constants.RedisCacheNullTTL)
		}
		return nil, err
	}

	// ========== record exist, save to redis and return ==========
	ttl := time.Minute * constants.RedisCacheShopTTL
	bytes, err := json.Marshal(shop)
	if err != nil {
		return nil, err
	}
	if err := s.rdb.Set(ctx, key, bytes, ttl).Err(); err != nil {
		return nil, err
	}
	if err := copier.CopyWithOption(&data, &shop, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *shopService) UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error {
	// ========== update sql db ==========
	// already check id with gin.Context shouldBind in handler
	var shop model.Shop
	if err := copier.CopyWithOption(&shop, &req, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	if err := s.shopRepository.Update(ctx, &shop); err != nil {
		return err
	}

	// ========== delete cache ==========
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.Id)
	s.rdb.Del(ctx, key)
	return nil
}
