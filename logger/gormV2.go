package logger

import (
	"context"
	"github.com/kataras/golog"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

func NewDBLogger(callback func(sql, file string, rows int64, msElapsed float64), leave gormLogger.LogLevel, slowThreshold time.Duration) gormLogger.Interface {
	return &dbLogger{
		LogLevel:      leave,
		SlowThreshold: slowThreshold,
		callback:      callback,
	}
}

type dbLogger struct {
	LogLevel      gormLogger.LogLevel
	SlowThreshold time.Duration
	callback      func(sql, file string, rows int64, msElapsed float64)
}

func (dbLogger *dbLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *dbLogger
	newLogger.LogLevel = level
	return &newLogger
}

func (dbLogger *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if dbLogger.LogLevel >= gormLogger.Info {
		golog.Default.Infof("Info: %s: %+v", msg, data)
	}
}

func (dbLogger *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if dbLogger.LogLevel >= gormLogger.Warn {
		golog.Default.Warnf("Warn: %s: %+v", msg, data)
	}
}

func (dbLogger *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if dbLogger.LogLevel >= gormLogger.Error {
		golog.Default.Errorf("Error: %s: %+v", msg, data)
	}
	return
}

const (
	DBFmtWithError   = "\033[36;1m[%.2fms]\033[0m \033[36;31m[%d rows affected or returned]\033[0m \n%s\n%s\n\033[31;1m%s\033[0m\n"
	DBFmtWithNoError = "\033[36;1m[%.2fms]\033[0m \033[36;31m[%d rows affected or returned]\033[0m \n%s\n%s\n"
)

func (dbLogger *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	file := utils.FileWithLineNum()
	sql, rows, msElapsed := dbLogger.getSqlInfo(fc, elapsed)
	switch {
	case err != nil && dbLogger.LogLevel >= gormLogger.Error:
		golog.Default.Errorf(DBFmtWithError, msElapsed, rows, sql, file, err.Error())
	case elapsed > dbLogger.SlowThreshold && dbLogger.SlowThreshold != 0 && dbLogger.LogLevel >= gormLogger.Warn:
		golog.Default.Warnf(DBFmtWithNoError, msElapsed, rows, sql, file)
	case dbLogger.LogLevel >= gormLogger.Info:
		golog.Default.Infof(DBFmtWithNoError, msElapsed, rows, sql, file)
	}
	if dbLogger.callback != nil {
		dbLogger.callback(sql, file, rows, msElapsed)
	}
}

func (dbLogger *dbLogger) getSqlInfo(fc func() (string, int64), elapsed time.Duration) (string, int64, float64) {
	sql, rows := fc()
	msElapsed := float64(elapsed.Nanoseconds()) / 1e6
	return sql, rows, msElapsed
}
