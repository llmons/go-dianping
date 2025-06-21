package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/entity"
	"gorm.io/gorm"
	"time"
)

type ShopService interface {
	GetShopById(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error)
	UpdateShop(ctx context.Context, req *v1.UpdateShopReq) error
}

func NewShopService(
	service *Service,
) ShopService {
	return &shopService{
		Service: service,
	}
}

type shopService struct {
	*Service
}

func (s *shopService) GetShopById(ctx context.Context, req *v1.GetShopByIDReq) (*v1.GetShopByIDRespData, error) {
	//data, err := s.queryWithPassThrough(ctx, req)
	//if err != nil {
	//	return nil, err
	//}

	data, err := s.queryWithMutex(ctx, req)
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

	if _, err := s.repo.Query.Shop.Where(s.repo.Query.Shop.ID.Eq(req.ID)).
		Updates(&shop); err != nil {
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
	shop, err := s.repo.Query.Shop.Where(s.repo.Query.Shop.ID.Eq(req.ID)).First()
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
	if isLock {
		time.Sleep(time.Millisecond * 50)
		return s.queryWithMutex(ctx, req)
	}

	// ========== query sql db ==========
	shop, err := s.repo.Query.Shop.Where(s.repo.Query.Shop.ID.Eq(req.ID)).First()
	// mock time delay
	time.Sleep(time.Millisecond * 200)
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
