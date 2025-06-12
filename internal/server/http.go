package server

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/handler"
	"go-dianping/internal/middleware"
	"go-dianping/pkg/log"
	"net/http"
)

func NewServerHTTP(
	logger *log.Logger,
	userHandler *handler.UserHandler,
	shopType *handler.ShopTypeHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(
		middleware.CORSMiddleware(),
	)
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	r.GET("/user", userHandler.GetUserById)

	r.GET("/shop-type/list", shopType.GetShopTypeList)
	return r
}
