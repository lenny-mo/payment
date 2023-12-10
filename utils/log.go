package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// 日志文件名称
	fileName := "go.micro.service.payment.log"
	syncWriter := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   fileName, // 文件名称
			MaxSize:    256,      // MB
			MaxBackups: 3,        // 最大备份
			LocalTime:  true,
			MaxAge:     30,
		})

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// core
	core := zapcore.NewCore(
		// 编码器
		zapcore.NewJSONEncoder(encoderConfig),
		syncWriter,
		// 设置日志级别
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	// 创建 logger 实例
	log := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	zap.ReplaceGlobals(log)
}
