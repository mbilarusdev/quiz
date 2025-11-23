package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Level.SetLevel(zap.DebugLevel)
	L, _ = cfg.Build()
}
