package logx

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logx struct {
	conf   *Conf
	logger *zap.Logger
}

func NewLogx(opts ...IOption) *Logx {
	conf := &Conf{}
	for _, v := range opts {
		v.apply(conf)
	}
	defaultConfig(conf)
	log := &Logx{conf: conf}
	log.logger = log.zap()
	return log
}

func (log *Logx) Zap() *zap.Logger {
	return log.logger
}

func (log *Logx) zap() *zap.Logger {
	opts := []zap.Option{
		zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel),
	}
	if len(log.conf.serviceName) > 0 {
		opts = append(opts, zap.Fields(
			zap.String("service", log.conf.serviceName),
		))
	}
	if log.conf.entry != nil {
		opts = append(opts, zap.Hooks(log.conf.entry))
	}
	var write zapcore.WriteSyncer
	switch log.conf.mode {
	case ModFile:
		write = log.sync()
	default:
		write = zapcore.Lock(os.Stdout)
	}
	return zap.New(zapcore.NewCore(log.encoder(), write, log.conf.level), opts...)
}

func (log *Logx) sync() zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", log.conf.path, log.conf.filename()),
		MaxSize:    log.conf.maxSize,
		MaxBackups: log.conf.maxBackups,
		Compress:   log.conf.compress,
		MaxAge:     log.conf.keepDay,
	})
}

func (log *Logx) encoder() zapcore.Encoder {
	var encoder zapcore.Encoder
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(log.conf.timeFormat))
	}
	switch log.conf.color {
	case true:
		conf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	default:
		conf.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	conf.TimeKey = timeKey
	switch log.conf.encoding {
	case EncodingJson:
		encoder = zapcore.NewJSONEncoder(conf)
	default:
		encoder = zapcore.NewConsoleEncoder(conf)
	}
	return encoder
}
