package redsync

import (
	"flag"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go-dianping/pkg/config"
	"os"
	"path/filepath"
	"testing"
)

var (
	conf *viper.Viper
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

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func TestRedSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})
	rdb2 := redis.NewClient(&redis.Options{
		Addr: conf.GetString("data.redis2.addr"),
		DB:   conf.GetInt("data.redis.db"),
	})
	rdb3 := redis.NewClient(&redis.Options{
		Addr: conf.GetString("data.redis3.addr"),
		DB:   conf.GetInt("data.redis.db"),
	})

	pool := goredis.NewPool(rdb)
	pool2 := goredis.NewPool(rdb2)
	pool3 := goredis.NewPool(rdb3)

	rs := redsync.New(pool, pool2, pool3)

	lock := rs.NewMutex("test")
	if err := lock.Lock(); err != nil {
		panic(err)
	}
	defer func(lock *redsync.Mutex) {
		if _, err := lock.Unlock(); err != nil {
			panic(err)
		}
	}(lock)
}
