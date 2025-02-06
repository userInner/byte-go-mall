package database

import (
	"byte-go-mall/common/logging"
	"byte-go-mall/constant/config"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

type DBLogger struct {
	logger        *logging.TracedLogger
	slowThreshold time.Duration
	config        *config.MySQLConfig
}

func newLogger(dbConfig *config.MySQLConfig) logger.Interface {

	return &DBLogger{
		logger:        logging.Logger,
		slowThreshold: time.Duration(dbConfig.SlowThreshold) * time.Millisecond,
	}
}

// LogMode 实现 logger.Interface
func (l *DBLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

// Info 实现 logger.Interface
func (l *DBLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(ctx, fmt.Sprintf(msg, data...))
}

// Warn 实现 logger.Interface
func (l *DBLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(ctx, fmt.Sprintf(msg, data...))
}

// Error 实现 logger.Interface
func (l *DBLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(ctx, fmt.Sprintf(msg, data...))
}

// Trace 实现 logger.Interface
func (l *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 记录慢查询
	if elapsed > l.slowThreshold {
		l.logger.Warn(ctx, "慢查询",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Error(err),
		)
		return
	}

	// 记录错误
	if err != nil {
		l.logger.Error(ctx, "查询错误",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Error(err),
		)
		return
	}

	// 记录正常查询
	if l.config.LogLevel == "debug" {
		l.logger.Debug(ctx, "SQL查询",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
	}
}
