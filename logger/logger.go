package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(debug bool) (*zap.Logger, error) {
	logLevel := zap.InfoLevel
	if debug {
		logLevel = zap.DebugLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	today := time.Now().Format("2006-01-02")
	outputFileName := "log/" + today + ".log"
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      debug,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout", outputFileName},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}

func NewLoggerForTest(debug bool) (*zap.Logger, error) {
	logLevel := zap.InfoLevel
	if debug {
		logLevel = zap.DebugLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      debug,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}
