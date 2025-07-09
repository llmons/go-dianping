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
	blogHandler *handler.BlogHandler,
	followHandler *handler.FollowHandler,
	shopHandler *handler.ShopHandler,
	shopTypeHandler *handler.ShopTypeHandler,
	uploadHandler *handler.UploadHandler,
	userHandler *handler.UserHandler,
	voucherHandler *handler.VoucherHandler,
	voucherOrderHandler *handler.VoucherOrderHandler,
) *http.Server {
	if conf.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.Default()
	g.RedirectTrailingSlash = true

	s := http.NewServer(
		g,
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
	{
		blogRouter := s.Group("/blog")
		{
			blogRouter.POST("", blogHandler.SaveBlog)
			blogRouter.PUT("/like/:id", blogHandler.LikeBlog)
			blogRouter.GET("/of/me", blogHandler.QueryMyBlog)
			blogRouter.GET("/hot", blogHandler.QueryHotBlog)
			blogRouter.GET("/:id", blogHandler.QueryById)
			blogRouter.GET("/of/follow", blogHandler.QueryBlogOfFollow)
		}

		followRouter := s.Group("/follow")
		{
			followRouter.PUT("/:id/:isFollow", followHandler.Follow)
			followRouter.GET("/or/not/:id", followHandler.IsFollow)
			followRouter.GET("/common/:id", followHandler.FollowCommons)
		}

		shopRouter := s.Group("/shop")
		{
			shopRouter.GET("/:id", shopHandler.QueryShopById)
			shopRouter.PUT("", shopHandler.UpdateShop)
		}

		shopTypeRouter := s.Group("/shop-type")
		{
			shopTypeRouter.GET("/list", shopTypeHandler.QueryTypeList)
		}

		uploadRouter := s.Group("/upload")
		{
			uploadRouter.POST("/blog", uploadHandler.UploadImage)
			uploadRouter.GET("/blog/delete", uploadHandler.DeleteBlogImg)
		}

		userRouter := s.Group("/user")
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

		voucherRouter := s.Group("/voucher")
		{
			voucherRouter.POST("/seckill", voucherHandler.AddSeckillVoucher)
			voucherRouter.POST("", voucherHandler.AddVoucher)
			voucherRouter.GET("/list/:shopId", voucherHandler.QueryVoucherOfShop)
		}

		voucherOrderRouter := s.Group("/voucher-order").Use(middleware.Login())
		{
			voucherOrderRouter.POST("/seckill/:id", voucherOrderHandler.SeckillVoucher)
		}
	}

	return s
}
