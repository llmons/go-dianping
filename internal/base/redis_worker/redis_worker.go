package redis_worker

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	beginTimestamp int64 = 1_640_995_200
	countBits      int64 = 32
)

type RedisWorker interface {
	NextId(ctx context.Context, keyPrefix string) (int64, error)
}

type redisWorker struct {
	rdb *redis.Client
}

func NewRedisWorker(rdb *redis.Client) RedisWorker {
	return &redisWorker{
		rdb: rdb,
	}
}

func (w *redisWorker) NextId(ctx context.Context, keyPrefix string) (int64, error) {
	//	 1. 生成时间戳
	timestamp := time.Now().Unix() - beginTimestamp

	//	2. 生成序列号
	//	2.1. 获取当天日期，精确到天
	date := time.Now().Format("2006-01-02")
	//	2.2. 自增长
	count, err := w.rdb.Incr(ctx, "icr:"+keyPrefix+":"+date).Result()
	if err != nil {
		return 0, err
	}

	//	3. 拼接并返回
	return timestamp<<countBits | count, nil
}
