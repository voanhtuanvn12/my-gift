package infra

import (
	"my-gift/configs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *configs.Config) (*zap.Logger, error) {
	var zapCfg zap.Config

	if cfg.App.Env == "production" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	switch cfg.Logger.Level {
	case "debug":
		zapCfg.Level.SetLevel(zapcore.DebugLevel)
	case "warn":
		zapCfg.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		zapCfg.Level.SetLevel(zapcore.ErrorLevel)
	default:
		zapCfg.Level.SetLevel(zapcore.InfoLevel)
	}

	return zapCfg.Build()
}
