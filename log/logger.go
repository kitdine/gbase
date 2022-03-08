package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger  Logger
	logLock sync.RWMutex
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

type Config struct {
	Level          string
	LogFileName    string
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

func init() {
	zapLoggerConfig := zap.NewDevelopmentConfig()
	zapLoggerEncoderConfig := zapcore.EncoderConfig{
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
	zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
	zapLogger, _ := zapLoggerConfig.Build(zap.AddCaller())
	SetLogger(&BaseLogger{zapLogger.Sugar()})
}

// InitLogger is init global logger
func InitLogger(config Config) (err error) {
	logLock.Lock()
	defer logLock.Unlock()
	logger, err = initLogger(config)
	return
}

func initLogger(config Config) (Logger, error) {
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
	return &BaseLogger{zapLogger.Sugar()}, nil
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

//SetLogger sets logger for sdk
func SetLogger(log Logger) {
	logLock.Lock()
	defer logLock.Unlock()
	logger = log
}

func GetLogger() Logger {
	logLock.RLock()
	defer logLock.RUnlock()
	return logger
}
