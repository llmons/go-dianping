package repository

import (
	"github.com/spf13/viper"
	"go-dianping/internal/query"
	"go-dianping/pkg/log"
	"go-dianping/pkg/zapgorm2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Repository struct {
	query  *query.Query
	logger *log.Logger
}

func NewRepository(query *query.Query, logger *log.Logger) *Repository {
	return &Repository{
		query:  query,
		logger: logger,
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
