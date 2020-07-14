// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package store

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type customLogger struct {
	log *zap.Logger
}

func newLogger() logger.Interface {
	return &customLogger{
		log: zap.L().WithOptions(zap.AddCallerSkip(2)).Named("store/sql"),
	}
}

func (l *customLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l customLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, data...))
}

func (l customLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, data...))
}

func (l customLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, data...))
}

func (l customLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
}
