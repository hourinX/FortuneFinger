package gormc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GORM日志适配器
type gormLogger struct {
	log           *zap.Logger
	level         gormlogger.LogLevel
	slowThreshold time.Duration
}

// 构建GORM日志器
func newGormLogger(log *zap.Logger, slowThreshold time.Duration) *gormLogger {
	return &gormLogger{
		log:           log,
		level:         gormlogger.Info,
		slowThreshold: slowThreshold,
	}
}

// 返回副本，各自持有日志级别
func (gl *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *gl
	clone.level = level
	return &clone
}

// 普通信息
func (gl *gormLogger) Info(_ context.Context, msg string, data ...any) {
	if gl.level >= gormlogger.Info {
		gl.log.Info(fmt.Sprintf(msg, data...))
	}
}

// 警告
func (gl *gormLogger) Warn(_ context.Context, msg string, data ...any) {
	if gl.level >= gormlogger.Warn {
		gl.log.Warn(fmt.Sprintf(msg, data...))
	}
}

// 错误
func (gl *gormLogger) Error(_ context.Context, msg string, data ...any) {
	if gl.level >= gormlogger.Error {
		gl.log.Error(fmt.Sprintf(msg, data...))
	}
}

// 记录SQL耗时、错误与慢查询
func (gl *gormLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if gl.level == gormlogger.Silent {
		return
	}
	duration := time.Since(begin)

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && gl.level >= gormlogger.Error:
		sql, rows := fc()
		gl.log.Error(
			"SQL ERROR",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Int64("dur_ms", duration.Milliseconds()),
			zap.Error(err),
		)
	case duration > gl.slowThreshold && gl.level >= gormlogger.Warn:
		sql, rows := fc()
		gl.log.Warn(
			"SQL SLOW",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Int64("dur_ms", duration.Milliseconds()),
		)
	case gl.level >= gormlogger.Info:
		sql, rows := fc()
		gl.log.Debug(
			"SQL INFO",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Int64("dur_ms", duration.Milliseconds()),
		)
	}
}
