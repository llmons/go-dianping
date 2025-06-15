//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-dianping/internal/handler"
	"go-dianping/internal/pkg/redis"
	"go-dianping/internal/repository"
	"go-dianping/internal/server"
	"go-dianping/internal/service"
	"go-dianping/pkg/log"
)

var ServerSet = wire.NewSet(
	redis.NewRedis,
	server.NewHttpServer,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewUserRepository,
	repository.NewShopTypeRepository,
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewShopTypeService,
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewShopTypeHandler,
)

func NewWire(*viper.Viper, *log.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		ServerSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
	))
}
