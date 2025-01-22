package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局logger实例
var log *zap.Logger

// 初始化logger配置
func init() {
	// 创建生产环境的logger配置
	config := zap.NewProductionConfig()

	// 自定义时间格式
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 构建logger
	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// Info 记录info级别的日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Error 记录error级别的日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 记录fatal级别的日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// With 创建带有预设字段的logger
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}
