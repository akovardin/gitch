package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level    string
	Encoding string
	Tags     []string
	Colored  bool
}

type Logger = zap.Logger

func New(config Config) *Logger {
	cfg := zap.NewProductionConfig()
	var lvl zapcore.Level
	err := lvl.UnmarshalText([]byte(config.Level))
	if err != nil {
		lvl = zap.InfoLevel
	}

	cfg.Level.SetLevel(lvl)
	cfg.DisableStacktrace = true
	cfg.Sampling.Initial = 50
	cfg.Sampling.Thereafter = 50
	cfg.Encoding = config.Encoding
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		return nil
	}))

	return logger.With(
		zap.Strings("tags", config.Tags),
	)
}
