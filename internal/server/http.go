package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	r.Use(
		middleware.CORSMiddleware(),
		sessions.Sessions("hmdp", store),
	)

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	api := r.Group("/api")
	{
		shopTypeRouter := api.Group("/shop-type")
		{
			shopTypeRouter.GET("/list", shopTypeHandler.GetShopTypeList)
		}

		userRouter := api.Group("/user")
		{
			userRouter.POST("/code", userHandler.SendCode)
		}
	}
	return r
}
