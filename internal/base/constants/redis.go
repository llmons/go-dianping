package constants

type RedisKey string
type RedisTTL int

// TTL minute
const (
	RedisLoginCodeKey = "login:code:"
	RedisLoginCodeTTL = 2

	RedisLoginUserKey = "login:token:"
	RedisLoginUserTTL = 30

	RedisCacheShopKey = "cache:shop:"
	RedisCacheShopTTL = 30

	RedisCacheShopTypeKey = "cache:shop-type"
	RedisCacheShopTypeTTL = 30

	RedisCacheNullTTL = 2
)
