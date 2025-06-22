package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/redis_data"
	"go-dianping/internal/entity"
	"go-dianping/internal/repository"
	"gorm.io/gorm"
	"time"
)

type ShopService interface {
	GetShopById(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error)
	UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error
	SaveShop2Redis(ctx context.Context, id int64, expireTime time.Duration) error
}

type shopService struct {
	*Service
	shopRepository   repository.ShopRepository
	cacheRebuildPool *ants.Pool
}

func NewShopService(
	service *Service,
	shopRepository repository.ShopRepository,
) ShopService {
	pool, err := ants.NewPool(10)
	if err != nil {
		return nil
	}

	return &shopService{
		Service:          service,
		shopRepository:   shopRepository,
		cacheRebuildPool: pool,
	}
}

func (s *shopService) GetShopById(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error) {
	// ========== solve Cache Penetration ==========
	//data, err := s.queryWithPassThrough(ctx, req)
	//if err != nil {
	//	return nil, err
	//}

	// ========== solve Hotspot Invalid with mutex ==========
	//data, err := s.queryWithMutex(ctx, req)
	//if err != nil {
	//	return nil, err
	//}

	// ========== solve Hotspot Invalid with logical expire ==========
	data, err := s.queryWithLogicalExpire(ctx, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *shopService) UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error {
	// ========== update sql db ==========
	// already check id with gin.Context shouldBind in handler
	var shop entity.Shop
	if err := copier.CopyWithOption(&shop, &req, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	if _, err := s.shopRepository.Update(ctx, &shop); err != nil {
		return err
	}

	// ========== delete cache ==========
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.ID)
	s.rdb.Del(ctx, key)
	return nil
}

func (s *shopService) queryWithPassThrough(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error) {
	var data v1.GetShopByIDRespData

	// ========== check cache ==========
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.ID)
	cacheShopStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// err == nil || err == redis.Nil
	if cacheShopStr != "" {
		var cacheShop entity.Shop
		if err := json.Unmarshal([]byte(cacheShopStr), &cacheShop); err != nil {
			return nil, err
		}
		if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
			return nil, err
		}
		return &data, nil
	}
	// cache str == ""
	if err == nil {
		return nil, errors.New("shop not exist")
	}

	// ========== query sql db ==========
	shop, err := s.shopRepository.GetById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// solve the cache penetration
			_, err := s.rdb.Set(ctx, key, "", constants.RedisCacheNullTTL).Result()
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	// ========== record exist, save to redis and return ==========
	bytes, err := json.Marshal(shop)
	if err != nil {
		return nil, err
	}
	_, err = s.rdb.Set(ctx, key, bytes, constants.RedisCacheShopTTL).Result()
	if err != nil {
		return nil, err
	}
	if err := copier.CopyWithOption(&data, &shop, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *shopService) queryWithMutex(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error) {
	var data v1.GetShopByIDRespData

	// ========== check cache ==========
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.ID)
	cacheShopStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// err == nil || err == redis.Nil
	if err == nil && cacheShopStr == "" {
		return nil, errors.New("shop not exist")
	}
	if cacheShopStr != "" {
		var cacheShop entity.Shop
		if err := json.Unmarshal([]byte(cacheShopStr), &cacheShop); err != nil {
			return nil, err
		}
		if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
			return nil, err
		}
		return &data, nil
	}

	// ========== lock ==========
	lockKey := fmt.Sprintf("%s%d", constants.RedisLockShopKey, req.ID)
	isLock, err := s.tryLock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	defer s.unlock(ctx, lockKey) // finally
	if !isLock {
		time.Sleep(time.Millisecond * 50)
		return s.queryWithMutex(ctx, req)
	}

	// ========== double check ==========
	cacheShopStr, err = s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// err == nil || err == redis.Nil
	if err == nil && cacheShopStr == "" {
		return nil, errors.New("shop not exist")
	}
	if cacheShopStr != "" {
		var cacheShop entity.Shop
		if err := json.Unmarshal([]byte(cacheShopStr), &cacheShop); err != nil {
			return nil, err
		}
		if err := copier.CopyWithOption(&data, &cacheShop, copier.Option{IgnoreEmpty: true}); err != nil {
			return nil, err
		}
		return &data, nil
	}

	// ========== query sql db ==========
	shop, err := s.shopRepository.GetById(ctx, req.ID)
	// mock time delay
	time.Sleep(time.Millisecond * 200)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// solve the cache penetration
			_, err := s.rdb.Set(ctx, key, "", constants.RedisCacheNullTTL).Result()
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	// ========== record exist, save to redis and return ==========
	bytes, err := json.Marshal(shop)
	if err != nil {
		return nil, err
	}
	_, err = s.rdb.Set(ctx, key, bytes, constants.RedisCacheShopTTL).Result()
	if err != nil {
		return nil, err
	}
	if err := copier.CopyWithOption(&data, &shop, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *shopService) queryWithLogicalExpire(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error) {
	var shopData v1.GetShopByIDRespData

	// ========== check cache ==========
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, req.ID)
	cacheDataStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// err == nil || err == redis.Nil
	if errors.Is(err, redis.Nil) || cacheDataStr == "" {
		return nil, errors.New("shop not exist")
	}

	var cacheData redis_data.RedisData[*entity.Shop]
	if err := json.Unmarshal([]byte(cacheDataStr), &cacheData); err != nil {
		return nil, err
	}
	if err := copier.Copy(&shopData, &cacheData.Data); err != nil {
		return nil, err
	}

	// ========== check expire time ==========\
	if cacheData.ExpireTime.After(time.Now()) {
		return &shopData, nil
	}

	// ========== lock ==========
	lockKey := fmt.Sprintf("%s%d", constants.RedisLockShopKey, req.ID)
	isLock, err := s.tryLock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if isLock {
		// ========== double check ==========
		cacheDataStr, err := s.rdb.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, err
		}
		// err == nil || err == redis.Nil
		if errors.Is(err, redis.Nil) || cacheDataStr == "" {
			return nil, errors.New("shop not exist")
		}

		// ========== rebuild cache ==========
		if err := s.cacheRebuildPool.Submit(func() {
			defer s.unlock(ctx, lockKey)
			newCtx := context.Background()
			if err := s.SaveShop2Redis(newCtx, shopData.ID, time.Second*20); err != nil {
				return
			}
		}); err != nil {
			return nil, err
		}
		return &shopData, nil
	}

	return &shopData, nil
}

func (s *shopService) tryLock(ctx context.Context, key string) (bool, error) {
	flag, err := s.rdb.SetNX(ctx, key, "1", time.Second*10).Result()
	if err != nil {
		return false, err
	}
	return flag, nil
}

func (s *shopService) unlock(ctx context.Context, key string) {
	s.rdb.Del(ctx, key)
}

func (s *shopService) SaveShop2Redis(ctx context.Context, id int64, expireTime time.Duration) error {
	shop, err := s.shopRepository.GetById(ctx, id)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 200)

	data := redis_data.RedisData[*entity.Shop]{
		ExpireTime: time.Now().Add(expireTime),
		Data:       shop,
	}

	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, shop.ID)
	_, err = s.rdb.Set(ctx, key, data, redis.KeepTTL).Result()
	if err != nil {
		return err
	}
	return nil
}
