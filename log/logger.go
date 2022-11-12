package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logLock    sync.RWMutex
	loggers    = make(map[string]Logger)
	global     Logger
	initialize bool
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

type GlobalConfig struct {
	Level       string
	LogFileName string
	ConfigBase
}

type ChildConfig struct {
	LoggerName  string
	Level       string
	LogFileName string
	ConfigBase
}

type ConfigBase struct {
	MaxSize        int
	MaxBackups     int
	MaxAge         int
	Compress       bool
	JSONFormat     bool
	ShowLineNumber bool
}

type BaseLogger struct {
	Logger
}

// Logger is the interface for Logger types
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
}

// InitLogger is init global logger
func InitLogger(config GlobalConfig, childConfig ...ChildConfig) {
	logLock.Lock()
	defer logLock.Unlock()
	initLogger(config)
	initialize = true
	initChildLoggers(childConfig...)
}

func AddChildLogger(config ...ChildConfig) {
	logLock.Lock()
	defer logLock.Unlock()
	if initialize {
		initChildLoggers(config...)
		return
	}
	panic("your should init global logger first! Please call InitLogger() first")
}

func initLogger(config GlobalConfig) {
	logLevel := getLogLevel(config.Level)
	encoderConfig := getEncoderConfig()
	var encoder zapcore.Encoder
	if config.JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	hook := getHooks(config.LogFileName, config.MaxSize, config.MaxBackups, config.MaxAge, config.Compress)

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), logLevel)
	zapLogger := zap.New(core)
	if config.ShowLineNumber {
		zapLogger = zapLogger.WithOptions(zap.AddCaller())
	}
	zap.ReplaceGlobals(zapLogger)
	global = &BaseLogger{zapLogger.Sugar()}
}

func initChildLoggers(childs ...ChildConfig) {
	for _, config := range childs {
		logLevel := getLogLevel(config.Level)
		encoderConfig := getEncoderConfig()
		var encoder zapcore.Encoder
		if config.JSONFormat {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}
		hook := getHooks(config.LogFileName, config.MaxSize, config.MaxBackups, config.MaxAge, config.Compress)
		childCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), logLevel)
		child := zap.L().With(zap.String("childLogName", config.LoggerName))
		child = child.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return childCore
		}))
		loggers[config.LoggerName] = &BaseLogger{child.Sugar()}
	}
}

func getLogLevel(level string) zapcore.Level {
	if zapLevel, ok := levelMap[level]; ok {
		return zapLevel
	}
	return zapcore.InfoLevel
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func getHooks(fileName string, maxSize, maxBackups, maxAge int, compress bool) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   fileName,   // 日志文件路径
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // 最多保留300个备份
		Compress:   compress,   // 是否压缩 disabled by default
		MaxAge:     maxAge,
	}
}

func GetLogger() Logger {
	logLock.RLock()
	defer logLock.RUnlock()
	return global
}

func GetLoggerWithName(name string) Logger {
	logLock.RLock()
	defer logLock.RUnlock()
	return loggers[name]
}