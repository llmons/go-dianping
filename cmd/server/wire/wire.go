//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-dianping/internal/handler"
	"go-dianping/internal/repository"
	"go-dianping/internal/server"
	"go-dianping/internal/service"
	"go-dianping/pkg/log"
	"go-dianping/pkg/redis"
)

var ServerSet = wire.NewSet(
	redis.NewRedis,
	server.NewHttpServer,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewQuery,
	repository.NewRepository,
	repository.NewUserRepository,
	repository.NewShopRepository,
	repository.NewShopTypeRepository,
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewShopService,
	service.NewShopTypeService,
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewShopHandler,
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
