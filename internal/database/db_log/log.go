package db_log

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

type Log struct {
	gormLog.Writer
	gormLog.Config
	Logger                              *zap.Logger
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

type Options interface {
	apply(*Log)
}

type OptionFunc func(log *Log)

type logOutPut struct {
	logger *zap.Logger
}

func New(options ...Options) *Log {
	logClass := &Log{
		infoStr:      "%s\n[info] ",
		warnStr:      "%s\n[warn] ",
		errStr:       "%s\n[error] ",
		traceStr:     "%s\n[%.2fms] [rows:%v] %s",
		traceWarnStr: "%s %s\n[%.2fms] [rows:%v] %s",
		traceErrStr:  "%s %s\n[%.2fms] [rows:%v] %s",
	}
	for _, val := range options {
		val.apply(logClass)
	}
	defaultOption(logClass)
	return logClass
}

func defaultOption(log *Log) {
	if log.Writer == nil {
		log.Writer = &logOutPut{
			logger: log.Logger,
		}
	}
	if (log.Config == gormLog.Config{}) {
		log.Config = gormLog.Config{
			SlowThreshold: 5 * time.Second,
			LogLevel:      gormLog.Warn,
			Colorful:      false,
		}
	}
}

func (l logOutPut) Printf(strFormat string, args ...interface{}) {
	logRes := fmt.Sprintf(strFormat, args...)
	if strings.HasPrefix(strFormat, "[info]") || strings.HasPrefix(strFormat, "[traceStr]") {
		l.logger.Info("gorm-->" + logRes)
	} else if strings.HasPrefix(strFormat, "[error]") || strings.HasPrefix(strFormat, "[traceErr]") {
		l.logger.Error("gorm-->" + logRes)
	} else if strings.HasPrefix(strFormat, "[warn]") || strings.HasPrefix(strFormat, "[traceWarn]") {
		l.logger.Warn("gorm-->" + logRes)
	}
}

// LogMode db_log mode
func (l *Log) LogMode(level gormLog.LogLevel) gormLog.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l Log) Info(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l Log) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l Log) Error(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l Log) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLog.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLog.Error && (!errors.Is(err, gormLog.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
		} else {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLog.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
		} else {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gormLog.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-1", sql)
		} else {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func (f OptionFunc) apply(log *Log) {
	f(log)
}

func SetInfoStrFormat(format string) Options {
	return OptionFunc(func(log *Log) {
		log.infoStr = format
	})
}

func SetWarnStrFormat(format string) Options {
	return OptionFunc(func(log *Log) {
		log.warnStr = format
	})
}

func SetErrStrFormat(format string) Options {
	return OptionFunc(func(log *Log) {
		log.errStr = format
	})
}

func SetTraceStrFormat(format string) Options {
	return OptionFunc(func(log *Log) {
		log.traceStr = format
	})
}
func SetTracWarnStrFormat(format string) Options {
	return OptionFunc(func(log *Log) {
		log.traceWarnStr = format
	})
}

func SetTracErrStrFormat(format string) Options {
	return OptionFunc(func(log *Log) {
		log.traceErrStr = format
	})
}

func SetWriter(writer gormLog.Writer) Options {
	return OptionFunc(func(log *Log) {
		log.Writer = writer
	})
}

func SetConfig(config gormLog.Config) Options {
	return OptionFunc(func(log *Log) {
		log.Config = config
	})
}

func SetLogger(logger *zap.Logger) Options {
	return OptionFunc(func(log *Log) {
		log.Logger = logger
	})
}
