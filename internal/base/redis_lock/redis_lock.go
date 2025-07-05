package redis_lock

import (
	"context"
	"github.com/duke-git/lancet/v2/random"
	"github.com/redis/go-redis/v9"
	"os"
	"path/filepath"
	"time"
)

var KeyLockPrefix = "redis_lock:"

type ILock interface {
	TryLock(ctx context.Context, timeout time.Duration) (bool, error)
	Unlock(ctx context.Context) error
}

type SimpleRedisLock struct {
	name         string
	uuid         string
	unlockScript *redis.Script
	rdb          *redis.Client
}

func NewSimpleRedisLock(name string, rdb *redis.Client) ILock {
	// 生成一个标识
	uuid, err := random.UUIdV4()
	if err != nil {
		return nil
	}

	// 加载 lua 脚本
	workDir, err := os.Getwd()
	if err != nil {
		return nil
	}
	luaPath := filepath.Join(workDir, "internal", "scripts", "unlock.lua")
	bytes, err := os.ReadFile(luaPath)
	if err != nil {
		return nil
	}
	unlockScript := redis.NewScript(string(bytes))

	return &SimpleRedisLock{
		name:         name,
		uuid:         uuid,
		unlockScript: unlockScript,
		rdb:          rdb,
	}
}

func (l *SimpleRedisLock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
	// 获取锁
	return l.rdb.SetNX(ctx, KeyLockPrefix+l.name, l.uuid, timeout).Result()
}

func (l *SimpleRedisLock) Unlock(ctx context.Context) error {
	//	调用 lua 脚本
	return l.unlockScript.Run(ctx, l.rdb, []string{KeyLockPrefix + l.name}, l.uuid).Err()
}

//func (l *SimpleRedisLock) Unlock(ctx context.Context) error {
//	// 获取锁中的标识
//	id, err := l.rdb.Get(ctx, KeyLockPrefix+l.name).Result()
//	if err != nil {
//		return err
//	}
//	// 判断标识是否一致
//	if l.uuid == id {
//		// 释放锁
//		return l.rdb.Del(ctx, KeyLockPrefix+l.name).Err()
//	}
//	return nil
//}
