package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugaredLogger *zap.SugaredLogger

func InitLogger(level string) {
	var atomicLevel zap.AtomicLevel

	switch level {
	case "info":
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		atomicLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		atomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	err := os.MkdirAll("/app/logger", 0755)
	if err != nil {
		log.Fatalf("Не удалось создать директорию логов: %v", err)
	}

	file, err := os.OpenFile("logger/info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	fileWriter := zapcore.AddSync(file)
	core := zapcore.NewCore(consoleEncoder, fileWriter, atomicLevel)

	logger := zap.New(core, zap.AddCaller())
	SugaredLogger = logger.Sugar()
}
