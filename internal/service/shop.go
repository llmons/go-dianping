package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/panjf2000/ants/v2"
	"go-dianping/api/v1"
	"go-dianping/internal/base/cache_client"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/entity"
	"go-dianping/internal/repository"
)

type ShopService interface {
	QueryById(ctx context.Context, req *v1.QueryShopByIDReq) (*v1.QueryShopByIDRespData, error)
	UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error
}

type shopService struct {
	*Service
	shopRepo         repository.ShopRepository
	cacheClient      cache_client.CacheClient[entity.Shop]
	cacheRebuildPool *ants.Pool
}

func NewShopService(
	service *Service,
	shopRepository repository.ShopRepository,
	cacheClient cache_client.CacheClient[entity.Shop],
) ShopService {
	return &shopService{
		Service:     service,
		shopRepo:    shopRepository,
		cacheClient: cacheClient,
	}
}

func (s *shopService) QueryById(ctx context.Context, req *v1.QueryShopByIDReq) (*v1.QueryShopByIDRespData, error) {
	// 解决缓存穿透
	shop, err := s.cacheClient.QueryWithPassThrough(ctx, constants.RedisCacheShopKey, *req.ID, s.shopRepo.GetById, constants.RedisCacheShopTTL)
	if err != nil {
		return nil, err
	}

	// 互斥锁解决缓存击穿
	//shop, err := s.cacheClient.QueryWithMutex(ctx, constants.RedisCacheShopKey, *req.ID, s.shopRepo.GetById, constants.RedisCacheShopTTL)
	//if err != nil {
	//	return nil, err
	//}

	// 逻辑过期解决缓存击穿
	//shop, err := s.cacheClient.QueryWithLogicalExpire(ctx, constants.RedisCacheShopKey, *req.ID, s.shopRepo.GetById, time.Second*20)
	//if err != nil {
	//	return nil, err
	//}

	// 返回
	var data v1.QueryShopByIDRespData
	if err := copier.Copy(&data, shop); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *shopService) UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error {
	if req.ID == nil {
		return v1.ErrShopIDIsNull
	}
	// 1. 更新数据库
	var shop entity.Shop
	if err := copier.Copy(&shop, &req); err != nil {
		return err
	}

	if _, err := s.shopRepo.Updates(ctx, &shop); err != nil {
		return err
	}
	// 2. 删除缓存
	return s.rdb.Del(ctx, fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.ID)).Err()
}
