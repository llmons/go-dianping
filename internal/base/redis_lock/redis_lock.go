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
	name string
	uuid string
	rdb  *redis.Client
}

func NewSimpleRedisLock(name string, rdb *redis.Client) ILock {
	return &SimpleRedisLock{
		name: name,
		rdb:  rdb,
	}
}

func (l *SimpleRedisLock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
	// 生成一个标识
	uuid, err := random.UUIdV4()
	if err != nil {
		return false, err
	}
	l.uuid = uuid
	// 获取锁
	return l.rdb.SetNX(ctx, KeyLockPrefix+l.name, l.uuid, timeout).Result()
}

func (l *SimpleRedisLock) Unlock(ctx context.Context) error {
	// 获取锁中的标识
	id, err := l.rdb.Get(ctx, KeyLockPrefix+l.name).Result()
	if err != nil {
		return err
	}
	// 判断标识是否一致
	if l.uuid == id {
		// 释放锁
		return l.rdb.Del(ctx, KeyLockPrefix+l.name).Err()
	}
	return nil
}
