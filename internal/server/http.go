package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-dianping/internal/handler"
	"go-dianping/internal/middleware"
	"net/http"
)

func NewHttpServer(
	rdb *redis.Client,
	userHandler *handler.UserHandler,
	shopTypeHandler *handler.ShopTypeHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(
		middleware.CORSMiddleware(),
	)

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
			shopRouter.GET("/")
		}

		shopTypeRouter := api.Group("/shop-type")
		{
			shopTypeRouter.GET("/list", shopTypeHandler.GetShopTypeList)
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
			authRouter := userRouter.Group("/").Use(middleware.Auth(rdb))
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
