package log

import (
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// TOPIC for setting topic of log
	TOPIC = "default-log"
	// LogTag default log tag
	LogTag = "default-api"

	//DEBUG ...
	DEBUG = iota
	//WARNING ...
	WARNING
	//INFO ...
	INFO
	//ERROR ...
	ERROR
	//FATAL ...
	FATAL
	//PANIC
	PANIC
)

var (
	logger    *Logger
	zapLogger *zap.SugaredLogger
)

type Logger struct {
	ZapLogger *zap.SugaredLogger
}

//DefaultLogger ...
func InitLogger() *Logger {
	if logger == nil {
		sugaredLog, _ := initSugaredLogger()
		l := Logger{
			ZapLogger: sugaredLog,
		}
		logger = &l
	}
	return logger
}

func (l *Logger) log(level int, message ...interface{}) {
	switch level {
	case DEBUG:
		l.ZapLogger.Debug(message...)
	case INFO:
		l.ZapLogger.Info(message...)
	case WARNING:
		l.ZapLogger.Warn(message...)
	case ERROR:
		l.ZapLogger.Error(message...)
	case FATAL:
		l.ZapLogger.Fatal(message...)
	case PANIC:
		l.ZapLogger.Panic(message...)
	}
}

// I ...
func I(message ...interface{}) {
	l := InitLogger()
	l.log(INFO, message)
}

// D ...
func D(message ...interface{}) {
	l := InitLogger()
	l.log(DEBUG, message)
}

// W ...
func W(message ...interface{}) {
	l := InitLogger()
	l.log(WARNING, message)
}

// E ...
func E(message ...interface{}) {
	l := InitLogger()
	l.log(ERROR, message)
}

// F ...
func F(message interface{}, context string, scope string) {
	l := InitLogger()
	l.log(ERROR, message, context, scope)
}

// P ...
func P(message ...interface{}) {
	l := InitLogger()
	l.log(ERROR, message)
}

func initSugaredLogger() (*zap.SugaredLogger, error) {
	var (
		logg *zap.Logger
		err  error
	)

	cfg := zap.Config{
		Encoding:         "json",
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey: "caller",
			EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
				_, caller.File, caller.Line, _ = runtime.Caller(8)
				enc.AppendString(caller.FullPath())
			},
		},
	}

	// check environment
	if os.Getenv("ENVIRONMENT") != "production" {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.Development = true
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		cfg.Development = false
	}

	logg, err = cfg.Build()
	if err != nil {
		return nil, err
	}

	defer logg.Sync()

	return logg.Sugar(), nil
}
