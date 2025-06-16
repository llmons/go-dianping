package constants

type RedisKey string
type RedisTTL int

const (
	RedisLoginCodeKey = "login:code:"
	RedisLoginCodeTTL = 2 // minute
	RedisLoginUserKey = "login:token:"
	RedisLoginUserTTL = 30 // minute
)
