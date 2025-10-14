package logging

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"pharmacy-modernization-project-model/internal/platform/config"
)

type LoggerBundle struct {
	Base *zap.Logger
}

func NewLogger(cfg *config.Config) *LoggerBundle {
	var l *zap.Logger

	// If logging is disabled, use a no-op logger
	if !cfg.Logging.Enabled {
		l = zap.NewNop()
		return &LoggerBundle{Base: l}
	}

	// Determine log level
	lvl := zapcore.DebugLevel
	_ = lvl.UnmarshalText([]byte(cfg.Logging.Level))

	// Determine encoder (console vs json)
	var encoder zapcore.Encoder
	if cfg.App.Env == "prod" || cfg.Logging.Format == "json" {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	// Determine output destination(s)
	var cores []zapcore.Core

	switch cfg.Logging.Output {
	case "file":
		// Write to file only
		fileWriter := getFileWriter(cfg)
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, lvl))

	case "both":
		// Write to both console and file
		consoleWriter := zapcore.AddSync(os.Stdout)
		fileWriter := getFileWriter(cfg)
		cores = append(cores,
			zapcore.NewCore(encoder, consoleWriter, lvl),
			zapcore.NewCore(encoder, fileWriter, lvl),
		)

	default: // "console"
		// Write to console only (default behavior)
		consoleWriter := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(encoder, consoleWriter, lvl))
	}

	// Combine cores and create logger
	core := zapcore.NewTee(cores...)
	l = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &LoggerBundle{Base: l}
}

func getFileWriter(cfg *config.Config) zapcore.WriteSyncer {
	// Ensure log directory exists
	logDir := filepath.Dir(cfg.Logging.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// If we can't create the directory, fall back to stdout
		return zapcore.AddSync(os.Stdout)
	}

	// Use lumberjack for log rotation
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Logging.FilePath,
		MaxSize:    cfg.Logging.FileMaxSize,    // megabytes
		MaxBackups: cfg.Logging.FileMaxBackups, // number of backups
		MaxAge:     cfg.Logging.FileMaxAge,     // days
		Compress:   true,                       // compress old files
	}

	return zapcore.AddSync(lumberJackLogger)
}
