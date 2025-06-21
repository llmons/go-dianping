package service_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/mock/gomock"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go-dianping/internal/repository"
	"go-dianping/internal/service"
	"go-dianping/pkg/config"
	"go-dianping/pkg/log"
	"go-dianping/pkg/redis"
	"os"
	"testing"
	"time"
)

var (
	conf   *viper.Viper
	logger *log.Logger
	rdb    *goRedis.Client
)

func TestMain(m *testing.M) {
	fmt.Println("begin")

	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		panic(err)
	}

	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf = config.NewConfig(*envConf)

	logger = log.NewLog(conf)
	rdb = redis.NewRedis(conf)

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func TestShopService_SaveShop2Redis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := repository.NewDB(conf, logger)
	query := repository.NewQuery(db)
	repo := repository.NewRepository(query, logger)
	shopRepo := repository.NewShopRepository(repo)
	srv := service.NewService(logger, conf, rdb)
	shopService := service.NewShopService(srv, shopRepo)

	ctx := context.Background()
	err := shopService.SaveShop2Redis(ctx, 1, time.Second*10)

	assert.NoError(t, err)
}
