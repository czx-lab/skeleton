package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Options struct {
	Debug            bool
	Encoder          string
	Entry            func(zapcore.Entry) error
	RecordTimeFormat string
	Filename         string
	MaxSize          int
	MaxBackups       int
	MaxAge           int
	Compress         bool
}

type Logger struct {
	*zap.Logger
	Options
}

// New 实例化Logger
func New(opts ...Option) (logger *zap.Logger, err error) {
	LoggerClass := &Logger{}
	for _, opt := range opts {
		opt(&LoggerClass.Options)
	}
	defaultConfig(&LoggerClass.Options)
	if LoggerClass.Options.Debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return
		}
		return
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig = LoggerClass.encoderConfig(encoderConfig)
	encoder := LoggerClass.setEncoder(encoderConfig)
	writer := LoggerClass.setLogSync()
	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	return zap.New(zapCore, zap.AddCaller(), zap.Hooks(LoggerClass.Options.Entry), zap.AddStacktrace(zap.WarnLevel)), nil
}

func (l *Logger) setEncoder(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
	var encoder zapcore.Encoder
	switch l.Options.Encoder {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return encoder
}

func (l *Logger) encoderConfig(encoderConfig zapcore.EncoderConfig) zapcore.EncoderConfig {
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(l.Options.RecordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "created_at"
	return encoderConfig
}

func (l *Logger) setLogSync() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   l.Options.Filename,
		MaxSize:    l.Options.MaxSize,
		MaxBackups: l.Options.MaxBackups,
		MaxAge:     l.Options.MaxAge,
		Compress:   l.Options.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

type Option func(opts *Options)

func defaultConfig(opts *Options) {
	if opts.Encoder == "" {
		opts.Encoder = "console"
	}
	if opts.RecordTimeFormat == "" {
		opts.RecordTimeFormat = "2006-01-02 15:04:05"
	}
	if opts.Filename == "" {
		curPath, _ := os.Getwd()
		opts.Filename = curPath + "/storage/logs/system.log"
	}
	if opts.MaxAge == 0 {
		opts.MaxAge = 10
	}
	if opts.Entry == nil {
		opts.Entry = func(entry zapcore.Entry) error {
			return nil
		}
	}
}

func WithDebug(debug bool) Option {
	return func(opts *Options) {
		opts.Debug = debug
	}
}

func WithEntry(entry func(zapcore.Entry) error) Option {
	return func(opts *Options) {
		opts.Entry = entry
	}
}

func WithEncode(encode string) Option {
	return func(opts *Options) {
		opts.Encoder = encode
	}
}

func WithRecordTimeFormat(recordTimeFormat string) Option {
	return func(opts *Options) {
		opts.RecordTimeFormat = recordTimeFormat
	}
}

func WithFilename(filename string) Option {
	return func(opts *Options) {
		opts.Filename = filename
	}
}

func WithMaxSize(maxSize int) Option {
	return func(opts *Options) {
		opts.MaxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(opts *Options) {
		opts.MaxBackups = maxBackups
	}
}

func WithMaxAge(maxAge int) Option {
	return func(opts *Options) {
		opts.MaxAge = maxAge
	}
}

func WithCompress(compress bool) Option {
	return func(opts *Options) {
		opts.Compress = compress
	}
}
