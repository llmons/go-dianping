package redis_lock

import (
	"context"
	"github.com/duke-git/lancet/v2/random"
	"github.com/redis/go-redis/v9"
	"time"
)

var KeyLockPrefix = "redis_lock:"

type ILock interface {
	TryLock(ctx context.Context, timeout time.Duration) (bool, error)
	Unlock(ctx context.Context) error
}

type SimpleRedisLock struct {
	rdb  *redis.Client
	name string
}

func NewSimpleRedisLock(name string, rdb *redis.Client) ILock {
	return &SimpleRedisLock{
		name: name,
		rdb:  rdb,
	}
}

func (l *SimpleRedisLock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
	// 生成一个标识
	id, err := random.UUIdV4()
	if err != nil {
		return false, err
	}
	// 获取锁
	return l.rdb.SetNX(ctx, KeyLockPrefix+l.name, id, timeout).Result()
}

func (l *SimpleRedisLock) Unlock(ctx context.Context) error {
	// 释放锁
	return l.rdb.Del(ctx, KeyLockPrefix+l.name).Err()
}
