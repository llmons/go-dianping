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
	query  *query.Query
	rdb    *redis.Client
}

func NewService(logger *log.Logger, conf *viper.Viper, query *query.Query, rdb *redis.Client) *Service {
	return &Service{
		logger: logger,
		conf:   conf,
		query:  query,
		rdb:    rdb,
	}
}
