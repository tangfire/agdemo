package data

import (
	"agdemo/internal/conf"
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewGreeterRepo, NewArticleRepo)

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

// data/redis.go
func NewRedis(c *conf.Data, logger log.Logger) (*redis.Client, error) {

	// 转换 durationpb.Duration 为 time.Duration
	readTimeout := c.Redis.ReadTimeout.AsDuration()
	writeTimeout := c.Redis.WriteTimeout.AsDuration()

	rdb := redis.NewClient(&redis.Options{
		Network:      c.Redis.Network,
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	// 测试连接
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	log.NewHelper(logger).Info("redis connect success")
	return rdb, nil
}

// NewData 整合所有数据源
// data/data.go
func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")

		// 关闭数据库连接
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}

		// 关闭Redis连接
		if rdb != nil {
			_ = rdb.Close()
		}
	}

	return &Data{
		db:  db,
		rdb: rdb,
		log: log.NewHelper(logger),
	}, cleanup, nil
}
