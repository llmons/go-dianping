package redis_data

import (
	"encoding/json"
	"time"
)

type RedisData struct {
	ExpireTime time.Time
	Data       any
}

func (d *RedisData) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *RedisData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}
