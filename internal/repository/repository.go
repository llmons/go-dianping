package repository

import (
	"github.com/spf13/viper"
	"go-dianping/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func NewDb(conf *viper.Viper) *gorm.DB {
	db, err := gorm.Open(mysql.Open(conf.GetString("data.mysql.dsn")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
