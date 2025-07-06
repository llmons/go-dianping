//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-dianping/internal/base/cache_client"
	"go-dianping/internal/base/redis_worker"
	"go-dianping/internal/handler"
	"go-dianping/internal/server"
	"go-dianping/internal/service"
	"go-dianping/pkg/app"
	"go-dianping/pkg/log"
	"go-dianping/pkg/server/http"
)

var cacheClientSet = wire.NewSet(
	cache_client.NewCacheClientForShop,
)

var redisWorkerSet = wire.NewSet(
	redis_worker.NewRedisWorker,
)

var serviceSet = wire.NewSet(
	service.NewDB,
	service.NewQuery,
	service.NewRedis,
	service.NewRedSync,
	service.NewService,
	service.NewSeckillVoucherService,
	service.NewShopService,
	service.NewShopTypeService,
	service.NewUserService,
	service.NewVoucherService,
	service.NewVoucherOrderService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewShopHandler,
	handler.NewShopTypeHandler,
	handler.NewUserHandler,
	handler.NewVoucherHandler,
	handler.NewVoucherOrderHandler,
)

var serverSet = wire.NewSet(
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
		cacheClientSet,
		redisWorkerSet,
		serviceSet,
		handlerSet,
		newApp,
	))
}
