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
	"time"
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
	pool, err := ants.NewPool(10)
	if err != nil {
		return nil
	}

	return &shopService{
		Service:          service,
		shopRepo:         shopRepository,
		cacheClient:      cacheClient,
		cacheRebuildPool: pool,
	}
}

func (s *shopService) QueryById(ctx context.Context, req *v1.QueryShopByIDReq) (*v1.QueryShopByIDRespData, error) {
	// 解决缓存穿透
	//shop, err := s.cacheClient.QueryWithPassThrough(ctx, constants.RedisCacheShopKey, *req.ID, s.shopRepo.GetById, constants.RedisCacheShopTTL)
	//if err != nil {
	//	return nil, err
	//}

	// 互斥锁解决缓存击穿

	// 逻辑过期解决缓存击穿
	shop, err := s.cacheClient.QueryWithLogicalExpire(ctx, constants.RedisCacheShopKey, *req.ID, s.shopRepo.GetById, time.Second*20)
	if err != nil {
		return nil, err
	}

	// 返回
	var data = v1.QueryShopByIDRespData{}
	if err := copier.Copy(&data, shop); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *shopService) UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error {
	// ========== update sql db ==========
	// already check id with gin.Context shouldBind in handler
	var shop entity.Shop
	if err := copier.CopyWithOption(&shop, &req, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	if _, err := s.shopRepo.Update(ctx, &shop); err != nil {
		return err
	}

	// ========== delete cache ==========
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.ID)
	s.rdb.Del(ctx, key)
	return nil
}
