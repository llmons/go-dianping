package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-dianping/docs"
	"go-dianping/internal/handler"
	"go-dianping/internal/middleware"
	"go-dianping/pkg/log"
	"net/http"
)

func NewHttpServer(
	logger *log.Logger,
	rdb *redis.Client,
	userHandler *handler.UserHandler,
	shopHandler *handler.ShopHandler,
	shopTypeHandler *handler.ShopTypeHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()
	r := gin.Default()

	// ========== swagger doc ==========
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	// ========== middleware ==========
	r.Use(
		middleware.CORSMiddleware(),
		middleware.RequestLogMiddleware(logger),
		middleware.ResponseLogMiddleware(logger),
		middleware.RefreshToken(rdb),
	)

	// ========== router ==========
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	api := r.Group("/api")
	{
		blogRouter := api.Group("/blog")
		{
			blogRouter.GET("/hot")
		}

		shopRouter := api.Group("/shop")
		{
			shopRouter.GET("/:id", shopHandler.GetShopById)
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
			voucherRouter.GET("/")
		}
	}
	return r
}
