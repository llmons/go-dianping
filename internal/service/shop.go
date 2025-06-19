package service

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go-dianping/api"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
	"strconv"
	"time"
)

type ShopService interface {
	GetShopById(ctx context.Context, params *api.GetShopByIdReq) (*api.GetShopByIdRespData, error)
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
		return &api.GetShopByIdRespData{
			Id:       strconv.Itoa(int(cacheShop.Model.Id)),
			Name:     cacheShop.Name,
			TypeId:   strconv.Itoa(cacheShop.TypeId),
			Images:   cacheShop.Images,
			Area:     cacheShop.Area,
			Address:  cacheShop.Address,
			X:        cacheShop.X,
			Y:        cacheShop.Y,
			AvgPrice: cacheShop.AvgPrice,
			Sold:     cacheShop.Sold,
			Comments: cacheShop.Comments,
			Score:    cacheShop.Score,
		}, nil
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
	return &api.GetShopByIdRespData{
		Id:       strconv.Itoa(int(shop.Model.Id)),
		Name:     shop.Name,
		TypeId:   strconv.Itoa(shop.TypeId),
		Images:   shop.Images,
		Area:     shop.Area,
		Address:  shop.Address,
		X:        shop.X,
		Y:        shop.Y,
		AvgPrice: shop.AvgPrice,
		Sold:     shop.Sold,
		Comments: shop.Comments,
		Score:    shop.Score,
	}, nil
}
