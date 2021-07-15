package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func New(logger string) *zap.Logger {
	return Logger.With(getFields(zap.String("module", logger))...)
}

// Configure 配置日志模块
func init() {
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "time",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewTee(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapcore.InfoLevel))
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Debug(msg string, fields ...zapcore.Field) {
	Logger.Debug(msg, getFields(fields...)...)
}

func Info(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, getFields(fields...)...)
}

func Warn(msg string, fields ...zapcore.Field) {
	Logger.Warn(msg, getFields(fields...)...)
}

func Error(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, getFields(fields...)...)
}

func Panic(msg string, fields ...zapcore.Field) {
	Logger.Panic(msg, getFields(fields...)...)
}

func getFields(fields ...zapcore.Field) []zapcore.Field {
	var f []zapcore.Field
	if len(fields) > 0 {
		f = append(f, fields...)
	}
	return f
}
