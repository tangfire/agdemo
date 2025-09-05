package main

import (
	myLogrus "github.com/go-kratos/kratos/contrib/log/logrus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

//// setupLogrus 配置 logrus 分级输出
//func setupLogrus() log.Logger {
//	// 创建 logrus 实例
//	logrusLogger := logrus.New()
//
//	// 设置日志格式为 JSON（可选）
//	logrusLogger.SetFormatter(&logrus.JSONFormatter{
//		TimestampFormat: "2006-01-02 15:04:05",
//	})
//
//	// 创建 info.log 文件（记录 INFO 及以上级别）
//	infoFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		panic(err)
//	}
//
//	// 创建 error.log 文件（记录 ERROR 及以上级别）
//	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		panic(err)
//	}
//
//	// 设置分级输出
//	logrusLogger.AddHook(&writerHook{
//		Writer:    infoFile,
//		LogLevels: []logrus.Level{logrus.InfoLevel, logrus.WarnLevel},
//	})
//	logrusLogger.AddHook(&writerHook{
//		Writer:    errorFile,
//		LogLevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
//	})
//
//	// 输出到控制台（可选）
//	logrusLogger.SetOutput(os.Stdout)
//
//	// 返回 Kratos 兼容的 Logger
//	return myLogrus.NewLogger(logrusLogger)
//}
//
//// writerHook 实现 logrus.Hook 接口，用于分级写入文件
//type writerHook struct {
//	Writer    *os.File
//	LogLevels []logrus.Level
//}
//
//// Fire 写入日志到文件
//func (hook *writerHook) Fire(entry *logrus.Entry) error {
//	line, err := entry.String()
//	if err != nil {
//		return err
//	}
//	_, err = hook.Writer.Write([]byte(line + "\n")) // 每条日志换行
//	return err
//}
//
//// Levels 返回监听的日志级别
//func (hook *writerHook) Levels() []logrus.Level {
//	return hook.LogLevels
//}

// setupLogrus 配置 logrus 分级输出（带日志轮转）
func setupLogrus() log.Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 配置 info.log 的日志轮转
	infoLog := &lumberjack.Logger{
		Filename:   "info.log", // 日志文件名
		MaxSize:    100,        // 单个文件最大大小（MB）
		MaxBackups: 30,         // 保留的旧日志文件最大数量
		MaxAge:     7,          // 保留旧日志的最大天数
		Compress:   true,       // 是否压缩/归档旧日志
	}

	// 配置 error.log 的日志轮转
	errorLog := &lumberjack.Logger{
		Filename:   "error.log",
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	}

	// 设置分级输出
	logrusLogger.AddHook(&writerHook{
		Writer:    infoLog,
		LogLevels: []logrus.Level{logrus.InfoLevel, logrus.WarnLevel},
	})
	logrusLogger.AddHook(&writerHook{
		Writer:    errorLog,
		LogLevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})

	// 输出到控制台（可选）
	logrusLogger.SetOutput(os.Stdout)

	return myLogrus.NewLogger(logrusLogger)
}

// writerHook 实现 logrus.Hook 接口（兼容 lumberjack）
type writerHook struct {
	Writer    *lumberjack.Logger // 改用 lumberjack.Logger
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line + "\n"))
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}
