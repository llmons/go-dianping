package service_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/mock/gomock"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go-dianping/internal/base/cache_client"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/entity"
	"go-dianping/internal/repository"
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

func TestShopService_SaveShop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := repository.NewDB(conf, logger)
	query := repository.NewQuery(db)
	repo := repository.NewRepository(query, logger)
	shopRepo := repository.NewShopRepository(repo)
	cacheClient := cache_client.NewCacheClient[entity.Shop](rdb)

	ctx := context.Background()
	shop, err := shopRepo.GetById(ctx, 1)
	if err != nil {
		return
	}
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, 1)
	cacheClient.SetWithLogicExpire(ctx, key, shop, time.Second*10)

	assert.NoError(t, err)
}
