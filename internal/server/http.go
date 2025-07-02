package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-dianping/docs"
	"go-dianping/internal/handler"
	"go-dianping/internal/middleware"
	"go-dianping/pkg/log"
	"go-dianping/pkg/server/http"
	netHttp "net/http"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	rdb *redis.Client,
	shopHandler *handler.ShopHandler,
	shopTypeHandler *handler.ShopTypeHandler,
	userHandler *handler.UserHandler,
	voucherHandler *handler.VoucherHandler,
	voucherOrderHandler *handler.VoucherOrderHandler,
) *http.Server {
	if conf.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		middleware.RefreshToken(rdb),
	)

	// ========== router ==========
	s.GET("/", func(ctx *gin.Context) {
		ctx.String(netHttp.StatusOK, "OK")
	})
	api := s.Group("/api")
	{
		blogRouter := api.Group("/blog")
		{
			blogRouter.GET("/hot")
		}

		shopRouter := api.Group("/shop")
		{
			shopRouter.GET("/:id", shopHandler.QueryShopById)
			shopRouter.PUT("/", shopHandler.UpdateShop)
		}

		shopTypeRouter := api.Group("/shop-type")
		{
			shopTypeRouter.GET("/list", shopTypeHandler.QueryTypeList)
		}

		uploadRouter := api.Group("/upload")
		{
			uploadRouter.GET("/")
		}

		userRouter := api.Group("/user")
		{
			noAuthRouter := userRouter.Group("/")
			{
				noAuthRouter.POST("/code", userHandler.SendCode)
				noAuthRouter.POST("/login", userHandler.Login)
			}
			authRouter := userRouter.Group("/").Use(middleware.Login())
			{
				authRouter.GET("/me", userHandler.Me)
			}
		}

		voucherRouter := api.Group("/voucher")
		{
			voucherRouter.POST("/seckill", voucherHandler.AddSeckillVoucher)
		}

		voucherOrderRouter := api.Group("/voucher-order")
		{
			voucherOrderRouter.POST("/seckill/:id", voucherOrderHandler.SeckillVoucher)
		}
	}

	return s
}
