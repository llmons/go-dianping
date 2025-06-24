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
	QueryById(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error)
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

func (s *shopService) QueryById(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error) {
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
	var data = v1.GetShopByIDRespData{}
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

//func (s *shopService) queryWithMutex(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error) {
//	var data v1.GetShopByIDRespData
//
//	// ========== check cache ==========
//	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.ID)
//	cacheShopStr, err := s.rdb.Get(ctx, key).Result()
//	if err != nil && !errors.Is(err, redis.Nil) {
//		return nil, err
//	}
//	// err == nil || err == redis.Nil
//	if err == nil && cacheShopStr == "" {
//		return nil, errors.New("shop not exist")
//	}
//	if cacheShopStr != "" {
//		var cacheShop entity.Shop
//		if err := json.Unmarshal([]byte(cacheShopStr), &cacheShop); err != nil {
//			return nil, err
//		}
//		if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
//			return nil, err
//		}
//		return &data, nil
//	}
//
//	// ========== lock ==========
//	lockKey := fmt.Sprintf("%s%d", constants.RedisLockShopKey, req.ID)
//	isLock, err := s.tryLock(ctx, lockKey)
//	if err != nil {
//		return nil, err
//	}
//	defer s.unlock(ctx, lockKey) // finally
//	if !isLock {
//		time.Sleep(time.Millisecond * 50)
//		return s.queryWithMutex(ctx, req)
//	}
//
//	// ========== double check ==========
//	cacheShopStr, err = s.rdb.Get(ctx, key).Result()
//	if err != nil && !errors.Is(err, redis.Nil) {
//		return nil, err
//	}
//	// err == nil || err == redis.Nil
//	if err == nil && cacheShopStr == "" {
//		return nil, errors.New("shop not exist")
//	}
//	if cacheShopStr != "" {
//		var cacheShop entity.Shop
//		if err := json.Unmarshal([]byte(cacheShopStr), &cacheShop); err != nil {
//			return nil, err
//		}
//		if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
//			return nil, err
//		}
//		return &data, nil
//	}
//
//	// ========== query sql db ==========
//	shop, err := s.shopRepo.GetById(ctx, *req.ID)
//	// mock time delay
//	time.Sleep(time.Millisecond * 200)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			// solve the cache penetration
//			_, err := s.rdb.Set(ctx, key, "", constants.RedisCacheNullTTL).Result()
//			if err != nil {
//				return nil, err
//			}
//		}
//		return nil, err
//	}
//
//	// ========== record exist, save to redis and return ==========
//	bytes, err := json.Marshal(shop)
//	if err != nil {
//		return nil, err
//	}
//	_, err = s.rdb.Set(ctx, key, bytes, constants.RedisCacheShopTTL).Result()
//	if err != nil {
//		return nil, err
//	}
//	if err := copier.CopyWithOption(&data, &shop, copier.Option{IgnoreEmpty: true}); err != nil {
//		return nil, err
//	}
//	return &data, nil
//}
