package redis_data

import "time"

type RedisData struct {
	ExpireTime time.Duration
	Data       any
}
