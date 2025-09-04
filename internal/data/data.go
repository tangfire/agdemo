package data

import (
	"agdemo/internal/conf"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewGreeterRepo, NewArticleRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db  *gorm.DB
	rdb *redis.Client
	log *log.Helper
}

// NewData .
//func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
//	cleanup := func() {
//		log.NewHelper(logger).Info("closing the data resources")
//	}
//	return &Data{}, cleanup, nil
//}

// NewDB 初始化数据库连接
func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: NewGormLogger(logger), // 自定义日志(见下文)
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(int(c.Database.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(c.Database.MaxOpenConns))
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// NewData 整合所有数据源
func NewData(c *conf.Data, db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}

	return &Data{
		db:  db,
		log: log.NewHelper(logger),
	}, cleanup, nil
}
