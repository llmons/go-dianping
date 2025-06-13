package repository

import (
	"github.com/spf13/viper"
	"go-dianping/pkg/log"
	"go-dianping/pkg/zapgorm2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
	//rdb    *redis.Client
	logger *log.Logger
}

func NewRepository(logger *log.Logger, db *gorm.DB) *Repository {
	return &Repository{
		db: db,
		//rdb:    rdb,
		logger: logger,
	}
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	dsn := conf.GetString("data.mysql.dsn")
	logger := zapgorm2.New(l.Logger).LogMode(gormlogger.Info)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
	return db
}
