package server

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-dianping/internal/dto"
	"go-dianping/internal/handler"
	"go-dianping/internal/middleware"
	"net/http"
)

func NewHttpServer(
	conf *viper.Viper,
	userHandler *handler.UserHandler,
	shopTypeHandler *handler.ShopTypeHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	store := cookie.NewStore([]byte(conf.GetString("security.session.key")))
	gob.Register(&dto.User{}) // register user struct for sessions
	r.Use(
		middleware.CORSMiddleware(),
		sessions.Sessions("hmdp", store),
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
			authRouter := userRouter.Group("/").Use(middleware.Auth())
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
