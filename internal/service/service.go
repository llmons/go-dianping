package service

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go-dianping/internal/query"
	"go-dianping/pkg/log"
	"go-dianping/pkg/zapgorm2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Service struct {
	logger *log.Logger
	conf   *viper.Viper
	query  *query.Query
	rdb    *redis.Client
	rs     *redsync.Redsync
}

func NewService(logger *log.Logger, conf *viper.Viper, query *query.Query, rdb *redis.Client, rs *redsync.Redsync) *Service {
	return &Service{
		logger: logger,
		conf:   conf,
		query:  query,
		rdb:    rdb,
		rs:     rs,
	}
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	dsn := conf.GetString("data.db.hmdp.dsn")
	logger := zapgorm2.New(l.Logger).LogMode(gormlogger.Info)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func NewQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}

func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
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

func NewRedSync(rdb *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(rdb)
	return redsync.New(pool)
}
