package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"
	"go-dianping/internal/base/cache_client"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
	"strconv"
)

type ShopService interface {
	QueryById(ctx context.Context, req *v1.QueryShopByIDReq) (*model.Shop, error)
	UpdateShop(ctx context.Context, req *model.Shop) error
	QueryShopOfType(ctx context.Context, typeId uint64, current int, x *float64, y *float64) ([]*model.Shop, error)
}

type shopService struct {
	*Service
	cacheClient cache_client.CacheClient[model.Shop]
}

func (s *shopService) QueryShopOfType(ctx context.Context, typeId uint64, current int, x *float64, y *float64) ([]*model.Shop, error) {
	if x == nil || y == nil {
		page, _, err := s.query.Shop.Where(s.query.Shop.TypeID.Eq(typeId)).
			FindByPage(current, constants.DefaultPageSize)
		return page, err
	}
	from := (current - 1) * constants.DefaultPageSize
	end := current * constants.DefaultPageSize
	key := fmt.Sprintf("%s%d", constants.RedisShopGeoKey, typeId)
	results, err := s.rdb.GeoSearchLocation(ctx, key, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  *y,
			Latitude:   *x,
			Radius:     5000,
			RadiusUnit: "m",
			Count:      end,
		},
		WithCoord: false,
		WithDist:  true,
		WithHash:  false,
	}).Result()
	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if len(results) <= from {
		return nil, nil
	}
	results = results[from:]
	ids := make([]uint64, len(results))
	distanceMap := map[string]float64{}
	for i, r := range results {
		shopIDStr := r.Name
		id, err := strconv.Atoi(shopIDStr)
		if err != nil {
			return nil, err
		}
		ids[i] = uint64(id)
		distanceMap[shopIDStr] = r.Dist
	}

	shops, err := s.query.Shop.Where(s.query.Shop.ID.In(ids...)).Order(s.query.Shop.ID.Field(ids...)).Find()
	if err != nil {
		return nil, err
	}
	for _, shop := range shops {
		shop.Distance = distanceMap[strconv.FormatUint(shop.ID, 10)]
	}

	return shops, nil
}

func NewShopService(
	service *Service,
	cacheClient cache_client.CacheClient[model.Shop],
) ShopService {
	return &shopService{
		Service:     service,
		cacheClient: cacheClient,
	}
}

func (s *shopService) QueryById(ctx context.Context, req *v1.QueryShopByIDReq) (*model.Shop, error) {
	// 解决缓存穿透
	shop, err := s.cacheClient.QueryWithPassThrough(ctx, constants.RedisCacheShopKey, *req.ID, s.query.Shop.GetByID, constants.RedisCacheShopTTL)
	if err != nil {
		return nil, err
	}

	// 互斥锁解决缓存击穿
	//shop, err := s.cacheClient.QueryWithMutex(ctx, constants.RedisCacheShopKey, *req.VoucherID, s.shopRepo.GetById, constants.RedisCacheShopTTL)
	//if err != nil {
	//	return nil, err
	//}

	// 逻辑过期解决缓存击穿
	//shop, err := s.cacheClient.QueryWithLogicalExpire(ctx, constants.RedisCacheShopKey, *req.VoucherID, s.shopRepo.GetById, time.Second*20)
	//if err != nil {
	//	return nil, err
	//}

	// 返回
	return shop, nil
}

func (s *shopService) UpdateShop(ctx context.Context, shop *model.Shop) error {
	// 1. 更新数据库
	if _, err := s.query.Shop.Where(s.query.Shop.ID.Eq(shop.ID)).Updates(&shop); err != nil {
		return err
	}
	// 2. 删除缓存
	return s.rdb.Del(ctx, fmt.Sprintf("%s%d", constants.RedisCacheShopKey, shop.ID)).Err()
}
