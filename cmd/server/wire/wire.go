//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-dianping/internal/base/cache_client"
	"go-dianping/internal/handler"
	"go-dianping/internal/repository"
	"go-dianping/internal/server"
	"go-dianping/internal/service"
	"go-dianping/pkg/app"
	"go-dianping/pkg/log"
	"go-dianping/pkg/redis"
	"go-dianping/pkg/server/http"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewQuery,
	repository.NewRepository,
	repository.NewSeckillVoucherRepository,
	repository.NewShopRepository,
	repository.NewShopTypeRepository,
	repository.NewUserRepository,
	repository.NewVoucherRepository,
)

var cacheClientSet = wire.NewSet(
	cache_client.NewCacheClientForShop,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewSeckillVoucherService,
	service.NewShopService,
	service.NewShopTypeService,
	service.NewUserService,
	service.NewVoucherService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewShopHandler,
	handler.NewShopTypeHandler,
	handler.NewUserHandler,
	handler.NewVoucherHandler,
)

var serverSet = wire.NewSet(
	redis.NewRedis,
	server.NewHTTPServer,
)

// build App
func newApp(
	httpServer *http.Server,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer),
		app.WithName("go-dianping"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		serverSet,
		repositorySet,
		cacheClientSet,
		serviceSet,
		handlerSet,
		newApp,
	))
}
