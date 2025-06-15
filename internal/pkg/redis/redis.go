package redis

import (
	"context"
	"fmt"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

var client *goRedis.Client

func NewRedis(conf *viper.Viper) *goRedis.Client {
	rdb := goRedis.NewClient(&goRedis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}
