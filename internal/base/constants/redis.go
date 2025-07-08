package constants

import "time"

type RedisKey string
type RedisTTL int

const (
	RedisLoginCodeKey = "login:code:"
	RedisLoginCodeTTL = time.Minute * 2

	RedisLoginUserKey = "login:token:"
	RedisLoginUserTTL = time.Minute * 30

	RedisCacheShopKey = "cache:shop:"
	RedisCacheShopTTL = time.Minute * 30

	RedisCacheShopTypeKey = "cache:shop-type"
	RedisCacheShopTypeTTL = time.Minute * 30

	RedisCacheNullTTL = time.Minute * 2

	RedisLockShopKey = "redis_lock:shop:"

	RedisSeckillStockKey = "seckill:stock:"

	RedisBlogLikeKey = "blog:liked:"
)
