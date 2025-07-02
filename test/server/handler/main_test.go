package handler

import (
	"flag"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-dianping/internal/handler"
	"go-dianping/internal/middleware"
	"go-dianping/internal/repository"
	"go-dianping/internal/service"
	"go-dianping/pkg/config"
	"go-dianping/pkg/log"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var logger *log.Logger
var hdl *handler.Handler
var router *gin.Engine

var db *gorm.DB
var rdb *redis.Client
var repo *repository.Repository

var srv *service.Service

func TestMain(m *testing.M) {
	fmt.Println("begin")
	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		fmt.Println("Setenv error", err)
	}
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	// modify log directory
	logPath := filepath.Join("../../../", conf.GetString("log.log_file_name"))
	conf.Set("log.log_file_name", logPath)

	logger = log.NewLog(conf)
	hdl = handler.NewHandler(logger)

	db = repository.NewDB(conf, logger)
	query := repository.NewQuery(db)
	repo = repository.NewRepository(query, logger)

	srv = service.NewService(logger, conf, rdb)

	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func newHttpExcept(t *testing.T, router *gin.Engine) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(router),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			// httpexpect.NewDebugPrinter(t, true),
		},
	})
}
