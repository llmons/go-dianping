package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go-dianping/internal/query"
	"go-dianping/pkg/log"
)

type Service struct {
	logger *log.Logger
	conf   *viper.Viper
	rdb    *redis.Client
	query  *query.Query
}

func NewService(logger *log.Logger, conf *viper.Viper, rdb *redis.Client, query *query.Query) *Service {
	return &Service{
		logger: logger,
		conf:   conf,
		rdb:    rdb,
		query:  query,
	}
}
