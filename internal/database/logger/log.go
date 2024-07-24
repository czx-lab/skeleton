package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	gormLog "gorm.io/gorm/logger"
)

type Log struct {
	gormLog.Writer
	gormLog.Config
	Logger *zap.Logger
}

type Options interface {
	apply(*Log)
}

type OptionFunc func(log *Log)

type logOutPut struct {
	logger *zap.Logger
}

func New(options ...Options) gormLog.Interface {
	xlog := &Log{}
	for _, val := range options {
		val.apply(xlog)
	}
	defaultOption(xlog)
	return gormLog.New(xlog.Writer, xlog.Config)
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

func (l logOutPut) Printf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if strings.HasPrefix(format, "[info]") {
		l.logger.Info(msg)
	} else if strings.HasPrefix(format, "[error]") {
		l.logger.Error(msg)
	} else if strings.HasPrefix(format, "[warn]") {
		l.logger.Warn(msg)
	} else {
		l.logger.Info(msg)
	}
}

func (f OptionFunc) apply(log *Log) {
	f(log)
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
