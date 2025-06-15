package service

import (
	"github.com/redis/go-redis/v9"
	"go-dianping/pkg/log"
)

type Service struct {
	logger *log.Logger
	rdb    *redis.Client
}

func NewService(logger *log.Logger, rdb *redis.Client) *Service {
	return &Service{
		logger: logger,
		rdb:    rdb,
	}
}
