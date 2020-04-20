// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package logger sets up zap logging system.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new zap logger replacing its global instances.
func New(config *Config) (*zap.Logger, error) {
	var level zapcore.Level

	err := level.Set(config.Level)
	if err != nil {
		return nil, err
	}

	zapCfg := &zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          string(config.Encoding),
		OutputPaths:       []string{config.Path},
		ErrorOutputPaths:  []string{config.Path},
	}

	switch config.Encoding {
	case EncodingConsole:
		zapCfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	case EncodingJSON:
		zapCfg.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)

	return logger, nil
}
