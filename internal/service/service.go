package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go-dianping/internal/query"
	"go-dianping/pkg/log"
	"go-dianping/pkg/zapgorm2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Service struct {
	logger *log.Logger
	conf   *viper.Viper
	query  *query.Query
	rdb    *redis.Client
}

func NewService(logger *log.Logger, conf *viper.Viper, query *query.Query, rdb *redis.Client) *Service {
	return &Service{
		logger: logger,
		conf:   conf,
		query:  query,
		rdb:    rdb,
	}
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	dsn := conf.GetString("data.db.hmdp.dsn")
	logger := zapgorm2.New(l.Logger).LogMode(gormlogger.Info)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func NewQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}
