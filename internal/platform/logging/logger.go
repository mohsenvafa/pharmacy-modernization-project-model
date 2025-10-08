package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"pharmacy-modernization-project-model/internal/platform/config"
)

type LoggerBundle struct {
	Base *zap.Logger
}

func NewLogger(cfg *config.Config) *LoggerBundle {
	var l *zap.Logger
	if cfg.App.Env == "prod" || cfg.Logging.Format == "json" {
		l, _ = zap.NewProduction()
	} else {
		conf := zap.NewDevelopmentConfig()
		lvl := zapcore.DebugLevel
		_ = lvl.UnmarshalText([]byte(cfg.Logging.Level))
		conf.Level = zap.NewAtomicLevelAt(lvl)
		l, _ = conf.Build()
	}
	return &LoggerBundle{Base: l}
}
