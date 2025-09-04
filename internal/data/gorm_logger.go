package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger log.Logger
}

func NewGormLogger(l log.Logger) logger.Interface {
	return &GormLogger{logger: l}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l // 可以根据需要实现日志级别控制
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.NewHelper(l.logger).Info(msg, data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.NewHelper(l.logger).Warn(msg, data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.NewHelper(l.logger).Error(msg, data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	log.NewHelper(l.logger).Infof("[%.3fms] rows:%v %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}
