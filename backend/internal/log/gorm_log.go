package log

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	level         logger.LogLevel
	slowThreshold time.Duration
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		level:         logger.Info,
		slowThreshold: 200 * time.Millisecond,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	nl := *l
	nl.level = level
	return &nl
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Info {
		Info(msg, args...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Warn {
		Warn(msg, args...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Error {
		l.print(colorRed, "ERROR", fmt.Sprintf(msg, args...))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	sql, rows := fc()
	elapsed := time.Since(begin)

	switch {
	case err != nil && l.level >= logger.Error:
		Error("%-40s → %10v  %d rows  err: %v", sql, elapsed, rows, err)
	case l.slowThreshold != 0 && elapsed > l.slowThreshold && l.level >= logger.Warn:
		Warn("%-40s → %10v  %d rows  (threshold: %s)", sql, elapsed, rows, l.slowThreshold)
	case l.level >= logger.Info:
		Info("%-40s → %10v  %d rows", sql, elapsed, rows)
	}
}

func (l *GormLogger) print(color, level, msg string) {
	Info("%s[GORM] [%-5s]%s %s", color, level, colorReset, msg)
}
