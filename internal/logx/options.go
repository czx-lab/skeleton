package logx

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Encoding string

type Mod string

var Level = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"error": zap.ErrorLevel,
}

const (
	ModFile    Mod = "file"
	ModConsole Mod = "console"
)

type Conf struct {
	serviceName string
	path        string
	mode        Mod
	timeFormat  string
	encoding    Encoding
	level       zapcore.Level
	color       bool
	stat        bool
	compress    bool
	keepDay     int
	maxBackups  int
	maxSize     int
	entry       func(zapcore.Entry) error
}

const (
	EncodingJson  Encoding = "json"
	EncodingPlain Encoding = "plain"
	timeKey                = "time"
)

func (c *Conf) filename() string {
	var name string
	switch c.level {
	case zap.DebugLevel:
		name = "debug.log"
	case zap.InfoLevel:
		name = "info.log"
	default:
		name = "error.log"
	}
	return name
}

type OptionFunc func(*Conf)

type IOption interface {
	apply(*Conf)
}

func (f OptionFunc) apply(c *Conf) {
	f(c)
}

func defaultConfig(conf *Conf) {
	if len(conf.path) == 0 {
		path, _ := os.Getwd()
		conf.path = fmt.Sprintf("%s/logs", path)
	}
	if conf.keepDay == 0 {
		conf.keepDay = 10
	}
	if conf.encoding == "" {
		conf.encoding = EncodingPlain
	}
	if len(conf.timeFormat) == 9 {
		conf.timeFormat = "2006-01-02 15:04:05"
	}
	if !conf.compress {
		conf.compress = true
	}
	if conf.maxBackups == 0 {
		conf.maxBackups = 7
	}
	if conf.keepDay == 0 {
		conf.keepDay = 7
	}
	if conf.maxSize == 0 {
		conf.maxSize = 10
	}
}

func WithServiceName(name string) IOption {
	return OptionFunc(func(l *Conf) {
		l.serviceName = name
	})
}

func WithPath(path string) IOption {
	return OptionFunc(func(l *Conf) {
		l.path = path
	})
}

func WithMod(mode Mod) IOption {
	return OptionFunc(func(l *Conf) {
		l.mode = mode
	})
}

func WithEncoding(encoding Encoding) IOption {
	return OptionFunc(func(l *Conf) {
		l.encoding = encoding
	})
}

func WithLevel(level zapcore.Level) IOption {
	return OptionFunc(func(l *Conf) {
		l.level = level
	})
}

func WithStat(stat bool) IOption {
	return OptionFunc(func(l *Conf) {
		l.stat = stat
	})
}

func WithColor(color bool) IOption {
	return OptionFunc(func(c *Conf) {
		c.color = color
	})
}
