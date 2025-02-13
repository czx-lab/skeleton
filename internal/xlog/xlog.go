package xlog

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const timeKey = "time"

var levels = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"error": zap.ErrorLevel,
	"warn":  zap.WarnLevel,
	"panic": zap.PanicLevel,
	"fatal": zap.FatalLevel,
}

type LogConf struct {
	ServiceName string
	Path        string
	Mode        string
	Encoding    string
	TimeFormat  string
	Level       string
	Compress    bool
	KeepDays    int
}

type XLog struct {
	conf LogConf

	instance *zap.Logger
}

func New(conf LogConf) *XLog {
	defaultConf(&conf)
	logger := &XLog{
		conf: conf,
	}
	logger.instance = logger.Logger()
	return logger
}

func (x *XLog) Logger() *zap.Logger {
	if x.instance == nil {
		x.instance = instance(x.conf)
	}
	return x.instance
}

func instance(conf LogConf) *zap.Logger {
	opts := []zap.Option{
		zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel),
	}
	if len(conf.ServiceName) > 0 {
		opts = append(opts, zap.Fields(zap.String("service", conf.ServiceName)))
	}

	var write zapcore.WriteSyncer
	switch conf.Mode {
	case "file":
		write = sync(conf)
	default:
		write = zapcore.Lock(os.Stdout)
	}

	// logs level default : debug
	level, ok := levels[conf.Level]
	if !ok {
		level = zap.DebugLevel
	}
	return zap.New(zapcore.NewCore(encoder(conf), write, level), opts...)
}

func sync(conf LogConf) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename: fmt.Sprintf("%s/%s", conf.Path, fmt.Sprintf("%s.log", conf.Level)),
		Compress: conf.Compress,
		MaxAge:   conf.KeepDays,
	})
}

func encoder(conf LogConf) zapcore.Encoder {
	var encoder zapcore.Encoder
	econf := zap.NewProductionEncoderConfig()
	econf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(conf.TimeFormat))
	}
	if conf.Level == "debug" {
		econf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		econf.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	econf.TimeKey = timeKey
	switch conf.Encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(econf)
	default:
		encoder = zapcore.NewConsoleEncoder(econf)
	}
	return encoder
}

func defaultConf(conf *LogConf) {
	if len(conf.Path) == 0 {
		path, _ := os.Getwd()
		conf.Path = fmt.Sprintf("%s/logs", path)
	}

	if len(conf.Level) == 0 {
		conf.Level = "debug"
	}

	if len(conf.Encoding) == 0 {
		conf.Encoding = "console"
	}

	if len(conf.TimeFormat) == 0 {
		conf.TimeFormat = "2006-01-02 15:04:05"
	}

	if !conf.Compress {
		conf.Compress = true
	}
}
