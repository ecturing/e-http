// 创建一个名为logutil的包
package logutil

import (
	"ews/config"
	"os"

	"github.com/rs/zerolog"
)

var (
	// 全局日志记录器实例
	Logger zerolog.Logger
	log    *logConfig
)

type logConfig struct {
	Level      zerolog.Level `yaml:"level"`
	Target     string        `yaml:"target"`
	TimeFormat string        `yaml:"timeFormat"`
}

func init() {
	e := config.Config.Decode(&log)
	if e != nil {
		Logger.Err(e).Msg("日志配置文件解析错误")
	}
	// 初始化全局日志记录器
	logFile, err := os.OpenFile(log.Target, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err) // 如果无法打开或创建文件，则程序终止
	}
	Logger = zerolog.New(logFile).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = log.TimeFormat
	zerolog.SetGlobalLevel(log.Level) // 可以根据需要设置默认的日志级别
}
