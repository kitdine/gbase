package log

import (
	"context"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	logLock    sync.RWMutex
	loggers    = make(map[string]*Logger)
	global     *Logger
	initialize bool
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

type GlobalConfig struct {
	Level         string
	EnableFileLog bool
	ConfigBase
	FileLogConfig
}

type ChildConfig struct {
	LoggerName    string
	Level         string
	EnableFileLog bool
	ConfigBase
	FileLogConfig
}

type ConfigBase struct {
	JSONFormat     bool
	ShowLineNumber bool
}

type FileLogConfig struct {
	FileName   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type Logger struct {
	l *zap.Logger
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

	var core zapcore.Core
	if config.EnableFileLog {
		hook := getHooks(config.FileName, config.MaxSize, config.MaxBackups, config.MaxAge, config.Compress)
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), logLevel)
	} else {
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel)
	}

	zapLogger := zap.New(core)
	if config.ShowLineNumber {
		zapLogger = zapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	zap.ReplaceGlobals(zapLogger)
	global = &Logger{zapLogger}
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

		var childCore zapcore.Core
		if config.EnableFileLog {
			hook := getHooks(config.FileName, config.MaxSize, config.MaxBackups, config.MaxAge, config.Compress)
			childCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), logLevel)
		} else {
			childCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel)
		}

		child := zap.L().Named(config.LoggerName)
		child = child.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return childCore
		}))
		if config.ShowLineNumber {
			child = child.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
		}
		loggers[config.LoggerName] = &Logger{child}
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

// GetLogger 获取全局Logger
func GetLogger() *Logger {
	logLock.RLock()
	defer logLock.RUnlock()
	return global
}

// GetLoggerWithName 根据自定义时定义的日志文件名获取日志Logger
func GetLoggerWithFileName(name string) *Logger {
	logLock.RLock()
	defer logLock.RUnlock()
	return loggers[name]
}

const logCtxKey = "log_with_context"

func GetLoggerWitchCtx(ctx context.Context) *Logger {
	log, ok := ctx.Value(logCtxKey).(*Logger)
	if ok {
		return log
	}
	return global
}
