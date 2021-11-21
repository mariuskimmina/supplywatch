package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() *Logger {
	logcfg := zap.NewProductionConfig()
	logcfg.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.Stamp))
	})
	log, _ := logcfg.Build()
	defer log.Sync()
	sugar := log.Sugar()
	return &Logger{
		sugar,
	}
}
