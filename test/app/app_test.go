package app

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
	"go-dianping/internal/base/redis_worker"
	"go-dianping/internal/model"
	"go-dianping/internal/service"
	"go-dianping/pkg/config"
	"go-dianping/pkg/log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var (
	conf   *viper.Viper
	logger *log.Logger
	rdb    *goRedis.Client

	redisWorker redis_worker.RedisWorker
)

func TestMain(m *testing.M) {
	fmt.Println("begin")

	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(workdir, "..", "..", "config", "local.yml")

	if err := os.Setenv("APP_CONF", configPath); err != nil {
		panic(err)
	}

	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf = config.NewConfig(*envConf)

	logger = log.NewLog(conf)
	rdb = service.NewRedis(conf)

	redisWorker = redis_worker.NewRedisWorker(rdb)

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func TestIdWork(t *testing.T) {
	var wg sync.WaitGroup
	task := func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			id, err := redisWorker.NextId(context.Background(), "shop")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("id = ", id)
		}
	}
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go task()
	}
	begin := time.Now().UnixMilli()
	wg.Wait()
	end := time.Now().UnixMilli()
	fmt.Println("time = ", end-begin)
}

func TestSaveSHop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := service.NewDB(conf, logger)
	query := service.NewQuery(db)
	cacheClient := cache_client.NewCacheClient[model.Shop](rdb)

	ctx := context.Background()
	shop, err := query.Shop.GetByID(1)
	if err != nil {
		return
	}
	key := fmt.Sprintf("%s%d", constants.RedisCacheShopKey, 1)
	if err := cacheClient.SetWithLogicExpire(ctx, key, shop, time.Second*10); err != nil {
		return
	}

	assert.NoError(t, err)
}
