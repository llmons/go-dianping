package redis_data

import (
	"encoding/json"
	"time"
)

// RedisData use generic type
// "any" type will influence copier
type RedisData[T any] struct {
	ExpireTime time.Time
	Data       T
}

func (d RedisData[T]) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d RedisData[T]) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &d)
}
