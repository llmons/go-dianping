package cache_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/redis_data"
	"gorm.io/gorm"
	"time"
)

var rdb *redis.Client

func Set(ctx context.Context, key string, value any, expireTime time.Duration) {
	rdb.Set(ctx, key, value, expireTime)
}

func SetWithLogicExpire[T any](ctx context.Context, key string, value T, expireTime time.Duration) {
	// 设置逻辑过期
	redisData := redis_data.RedisData[T]{
		ExpireTime: time.Now().Add(expireTime),
		Data:       value,
	}
	// 写入 redis
	rdb.Set(ctx, key, redisData, redis.KeepTTL)
}

func QueryWithPassThrough[R any](ctx context.Context, keyPrefix string, id int64, dbFallback func(context.Context, int64) (*R, error), expireTime time.Duration) (*R, error) {
	key := fmt.Sprintf("%s%d", keyPrefix, id)
	// 1. 从 redis 查询店铺缓存
	_json, err := rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	// 2. 判断是否存在
	if _json != "" {
		// 3. 存在，直接返回
		var ret R
		if err := json.Unmarshal([]byte(_json), &ret); err != nil {
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
			rdb.Set(ctx, key, "", constants.RedisCacheNullTTL)
		}
		// 返回错误信息
		return nil, err
	}

	// 6. 存在，写入 redis
	Set(ctx, key, ret, expireTime)
	return ret, nil
}
