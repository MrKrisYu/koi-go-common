package entity

import (
	"github.com/MrKrisYu/koi-go-common/logger"
	"github.com/MrKrisYu/koi-go-common/logger/zap"
	"github.com/MrKrisYu/koi-go-common/sdk"
)

type Log struct {
	Driver     string `json:"driver"`     // 日志实现： default, zap
	Level      string `json:"level"`      // 日志过滤等级
	FilePath   string `json:"filePath"`   // 存放日志的路径
	MaxBackups int    `json:"maxBackups"` // 日志保留最大数量
	MaxAge     int    `json:"maxAge"`     // 日志保留最长时间，单位：天
	MaxSize    int    `json:"maxSize"`    // 日志保留的最大大小，单位：MB
	Compress   bool   `json:"compress"`   // 是否压缩日志
}

var LogConfig = new(Log)

const (
	DefaultLogDriver = "default"
	ZapLogDriver     = "zap"
)

func getDefaultLogConfig() *Log {
	return &Log{
		Driver:     "default",
		Level:      "info",
		FilePath:   "./logs/application.log",
		MaxBackups: 30,
		MaxAge:     30,
		MaxSize:    1024,
		Compress:   true,
	}
}

func (l *Log) Key() string {
	return "log"
}

func (l *Log) Setup() {
	defaultConfig := getDefaultLogConfig()
	if len(l.Driver) == 0 {
		l.Driver = defaultConfig.Driver
	}
	if len(l.Level) == 0 {
		l.Level = defaultConfig.Level
	}
	if len(l.FilePath) == 0 {
		l.FilePath = defaultConfig.FilePath
	}
	if l.MaxBackups <= 0 {
		l.MaxBackups = defaultConfig.MaxBackups
	}
	if l.MaxAge <= 0 {
		l.MaxAge = defaultConfig.MaxAge
	}
	if l.MaxSize <= 0 {
		l.MaxSize = defaultConfig.MaxSize
	}

	switch l.Driver {
	case DefaultLogDriver:
		defaultLogger := logger.NewDefaultLogger(logger.WithLevel(logger.InfoLevel))
		sdk.RuntimeContext.SetLogger(logger.NewHelper(defaultLogger))
		return
	case ZapLogDriver:
		newLogger, err := zap.NewLogger(
			logger.WithLevel(logger.InfoLevel),
			zap.WithOutput(
				zap.GetConsoleWriterSync(),
				zap.GetFileWriterSync(l.FilePath, l.MaxSize, l.MaxBackups, l.MaxAge, l.Compress)),
		)
		if err != nil {
			panic(err)
		}
		sdk.RuntimeContext.SetLogger(logger.NewHelper(newLogger))
		return
	}
}
