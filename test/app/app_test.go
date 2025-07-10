package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
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
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	conf   *viper.Viper
	logger *log.Logger
	rdb    *redis.Client

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

func TestLoadShopData(t *testing.T) {
	db := service.NewDB(conf, logger)
	query := service.NewQuery(db)

	// 1. 查询店铺信息
	list, err := query.Shop.Find()
	if err != nil {
		return
	}
	// 2. 将店铺信息存入缓存
	group := map[uint64][]*model.Shop{}
	for _, shop := range list {
		if _, ok := group[shop.ID]; !ok {
			group[shop.ID] = []*model.Shop{}
		}
		group[shop.TypeID] = append(group[shop.TypeID], shop)
	}

	ctx := context.Background()
	// 3. 将店铺类型信息存入缓存
	for typeId, value := range group {
		key := fmt.Sprintf("%s%d", constants.RedisShopGeoKey, typeId)
		var locations []*redis.GeoLocation

		for _, shop := range value {
			locations = append(locations, &redis.GeoLocation{
				Name:      strconv.Itoa(int(shop.ID)),
				Longitude: shop.Y,
				Latitude:  shop.X,
			})

		}
		rdb.GeoAdd(ctx, key, locations...)
	}

	assert.NoError(t, err)
}

func TestHyperLogLog(t *testing.T) {
	ctx := context.Background()
	values := make([]any, 1000)
	var err error
	for i := 0; i < 1_000_000; i++ {
		j := i % 1000
		values[j] = "user_" + strconv.Itoa(i)
		if j == 999 {
			err = rdb.PFAdd(ctx, "h12", values...).Err()
		}
	}

	assert.NoError(t, err)
}
