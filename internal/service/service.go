package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go-dianping/internal/repository"
	"go-dianping/pkg/log"
)

type Service struct {
	logger *log.Logger
	conf   *viper.Viper
	rdb    *redis.Client
	tm     repository.Transaction
}

func NewService(logger *log.Logger, conf *viper.Viper, rdb *redis.Client, tm repository.Transaction) *Service {
	return &Service{
		logger: logger,
		conf:   conf,
		rdb:    rdb,
		tm:     tm,
	}
}
