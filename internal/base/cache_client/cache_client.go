package cache_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/redis_data"
	"go-dianping/internal/entity"
	"gorm.io/gorm"
	"time"
)

type CacheClient[ENTITY any] interface {
	Set(ctx context.Context, key string, value any, expireTime time.Duration)
	SetWithLogicExpire(ctx context.Context, key string, value *ENTITY, expireTime time.Duration)
	QueryWithPassThrough(ctx context.Context, keyPrefix string, id int64, dbFallback func(context.Context, int64) (*ENTITY, error), expireTime time.Duration) (*ENTITY, error)
	QueryWithMutex(ctx context.Context, keyPrefix string, id int64, dbFallback func(context.Context, int64) (*ENTITY, error), expireTime time.Duration) (*ENTITY, error)
	QueryWithLogicalExpire(ctx context.Context, keyPrefix string, id int64, dbFallback func(context.Context, int64) (*ENTITY, error), expireTime time.Duration) (*ENTITY, error)
}
type cacheClient[ENTITY any] struct {
	rdb              *redis.Client
	cacheRebuildPool *ants.Pool
}

func NewCacheClient[ENTITY any](rdb *redis.Client) CacheClient[ENTITY] {
	pool, err := ants.NewPool(10)
	if err != nil {
		return nil
	}

	return &cacheClient[ENTITY]{
		rdb:              rdb,
		cacheRebuildPool: pool,
	}
}

func (c *cacheClient[ENTITY]) Set(ctx context.Context, key string, value any, expireTime time.Duration) {
	c.rdb.Set(ctx, key, value, expireTime)
}

func (c *cacheClient[ENTITY]) SetWithLogicExpire(ctx context.Context, key string, value *ENTITY, expireTime time.Duration) {
	// 设置逻辑过期
	redisData := redis_data.RedisData[*ENTITY]{
		ExpireTime: time.Now().Add(expireTime),
		Data:       value,
	}
	// 写入 redis
	c.rdb.Set(ctx, key, redisData, redis.KeepTTL)
}

func (c *cacheClient[ENTITY]) QueryWithPassThrough(ctx context.Context, keyPrefix string, id int64, dbFallback func(context.Context, int64) (*ENTITY, error), expireTime time.Duration) (*ENTITY, error) {
	key := fmt.Sprintf("%s%d", keyPrefix, id)
	// 1. 从 redis 查询店铺缓存
	shopJson, err := c.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// 2. 判断是否存在
	if shopJson != "" {
		// 3. 存在，直接返回
		var ret ENTITY
		if err := json.Unmarshal([]byte(shopJson), &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
	// 判断命中的是否是空值
	if err == nil {
		// 返回一个错误信息
		return nil, gorm.ErrRecordNotFound
	}

	// 4. 不存在，根据 id 查询数据库
	ret, err := dbFallback(ctx, id)
	if err != nil {
		// 5. 不存在，返回错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 将空值写入 redis
			c.rdb.Set(ctx, key, "", constants.RedisCacheNullTTL)
		}
		// 返回错误信息
		return nil, err
	}

	// 6. 存在，写入 redis
	c.Set(ctx, key, ret, expireTime)
	return ret, nil
}

func (c *cacheClient[ENTITY]) QueryWithMutex(ctx context.Context, keyPrefix string, id int64, dbFallback func(context.Context, int64) (*ENTITY, error), expireTime time.Duration) (*ENTITY, error) {
	key := fmt.Sprintf("%s%d", keyPrefix, id)
	// 1. 从 redis 查询店铺缓存
	shopJson, err := c.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// 2. 判断是否存在
	if shopJson != "" {
		// 3. 存在，直接返回
		var ret ENTITY
		if err := json.Unmarshal([]byte(shopJson), &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
	// 判断命中的是否是空值
	if err == nil {
		// 返回一个错误信息
		return nil, gorm.ErrRecordNotFound
	}

	// 4. 实现缓存重建
	// 4.1. 获取互斥锁
	lockKey := fmt.Sprintf("%s%d", constants.RedisLockShopKey, id)
	isLock, err := c.tryLock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	// 等待释放锁
	defer c.unlock(ctx, lockKey)
	if !isLock {
		// 4.3. 获取锁失败，休眠并重试
		time.Sleep(time.Millisecond * 50)
		return c.QueryWithMutex(ctx, keyPrefix, id, dbFallback, expireTime)
	}

	// ========== Double Check ==========
	// 1. 从 redis 查询店铺缓存
	shopJson, err = c.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// 2. 判断是否存在
	if shopJson != "" {
		// 3. 存在，直接返回
		var ret ENTITY
		if err := json.Unmarshal([]byte(shopJson), &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
	// 判断命中的是否是空值
	if err == nil {
		// 返回一个错误信息
		return nil, gorm.ErrRecordNotFound
	}

	// 4.4. 获取锁成功，根据 id 查询数据库
	shop, err := dbFallback(ctx, id)
	// 5. 不存在，返回错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 将空值写入 redis
			c.rdb.Set(ctx, key, "", constants.RedisCacheNullTTL)
		}
		// 返回错误信息
		return nil, err
	}

	// 6. 存在，写入 redis
	err = c.rdb.Set(ctx, key, shop, expireTime).Err()
	if err != nil {
		return nil, err
	}

	// 8. 返回
	return shop, nil
}

func (c *cacheClient[ENTITY]) QueryWithLogicalExpire(ctx context.Context, keyPrefix string, id int64, dbFallback func(context.Context, int64) (*ENTITY, error), expireTime time.Duration) (*ENTITY, error) {
	key := fmt.Sprintf("%s%d", keyPrefix, id)
	// 1. 从 redis 查询商铺缓存
	shopJson, err := c.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// 2. 判断是否存在
	if errors.Is(err, redis.Nil) || shopJson == "" {
		// 3. 不存在，返回错误
		return nil, errors.New("shop not exist")
	}

	// 4. 命中，需要先把 json 反序列化为对象
	var redisData redis_data.RedisData[ENTITY]
	if err := json.Unmarshal([]byte(shopJson), &redisData); err != nil {
		return nil, err
	}
	var ret ENTITY
	if err := copier.Copy(&ret, &redisData.Data); err != nil {
		return nil, err
	}

	// 5. 判断是否过期
	if redisData.ExpireTime.After(time.Now()) {
		// 5.1. 未过期，直接返回店铺信息
		return &ret, nil
	}
	// 5.2. 已过期，需要缓存重建
	// 6. 缓存重建
	// 6.1. 获取互斥锁
	lockKey := fmt.Sprintf("%s%d", constants.RedisLockShopKey, id)
	isLock, err := c.tryLock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	// 6.2. 判断是否获取锁成功
	if isLock {
		// ========== Double Check ==========
		shopJson, err := c.rdb.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, err
		}
		if errors.Is(err, redis.Nil) || shopJson == "" {
			return nil, errors.New("shop not exist")
		}

		// 6.3. 成功，开启独立线程，实现缓存重建
		if err := c.cacheRebuildPool.Submit(func() {
			// 等待释放锁
			defer c.unlock(ctx, lockKey)
			// 查询数据库
			newCtx := context.Background()
			retFromDB, err := dbFallback(newCtx, id)
			if err != nil {
				return
			}
			// 写入 redis
			c.SetWithLogicExpire(newCtx, key, retFromDB, expireTime)
		}); err != nil {
			return nil, err
		}
		// 6.4. 返回过期的商铺信息
		return &ret, nil
	}

	return &ret, nil
}

func (c *cacheClient[ENTITY]) tryLock(ctx context.Context, key string) (bool, error) {
	flag, err := c.rdb.SetNX(ctx, key, "1", time.Second*10).Result()
	if err != nil {
		return false, err
	}
	return flag, nil
}

func (c *cacheClient[ENTITY]) unlock(ctx context.Context, key string) {
	c.rdb.Del(ctx, key)
}

func NewCacheClientForShop(rdb *redis.Client) CacheClient[entity.Shop] {
	return NewCacheClient[entity.Shop](rdb)
}
