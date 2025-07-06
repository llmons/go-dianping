package service

import (
	"context"
	"fmt"
	"go-dianping/api/v1"
	"go-dianping/internal/base/cache_client"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
)

type ShopService interface {
	QueryById(ctx context.Context, req *v1.QueryShopByIDReq) (*model.Shop, error)
	UpdateShop(ctx context.Context, req *model.Shop) error
}

type shopService struct {
	*Service
	cacheClient cache_client.CacheClient[model.Shop]
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
