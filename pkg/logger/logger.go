package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var GlobalLogger *zap.Logger

var GlobalSugaredLogger *zap.SugaredLogger

func LoggerInit() {
	lgCfg := DefaultZapLoggerConfig
	logger, err := lgCfg.Build()
	if err != nil {
		panic(fmt.Errorf("init logger failed: %w", err))
	}
	zap.ReplaceGlobals(logger)
	GlobalLogger = logger
	GlobalSugaredLogger = logger.Sugar()
}

func Named(name string) Logger {
	return GlobalLogger.Named(name).Sugar()
}

var DefaultLevel = zap.InfoLevel

var DefaultZapLoggerConfig = zap.Config{
	Level:       zap.NewAtomicLevelAt(DefaultLevel),
	Development: true,
	Sampling: &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	},

	Encoding: "console",

	// copied from "zap.NewProductionEncoderConfig" with some updates
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,

		// Custom EncodeTime function to ensure we match format and precision of historic capnslog timestamps
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},

		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	},

	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}
