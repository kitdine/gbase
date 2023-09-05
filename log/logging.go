package log

import (
	"go.uber.org/zap"
)

func Named(s string) *Logger {
	l := &Logger{
		l: GetLogger().l,
	}
	return l.Named(s)
}

func With(field ...zap.Field) *Logger {
	l := &Logger{
		l: GetLogger().l,
	}
	return l.With(field...)
}

func Debug(msg string, field ...zap.Field) {
	GetLogger().Debug(msg, field...)
}

func Info(msg string, field ...zap.Field) {
	GetLogger().Info(msg, field...)
}

func Warn(msg string, field ...zap.Field) {
	GetLogger().Warn(msg, field...)
}

func Error(msg string, field ...zap.Field) {
	GetLogger().Error(msg, field...)
}

func DPanic(msg string, field ...zap.Field) {
	GetLogger().DPanic(msg, field...)
}

func Panic(msg string, field ...zap.Field) {
	GetLogger().Panic(msg, field...)
}

func Fatal(msg string, field ...zap.Field) {
	GetLogger().Fatal(msg, field...)
}

func (log *Logger) Named(s string) *Logger {
	log.l = log.l.Named(s)
	return log
}

func (log *Logger) With(field ...zap.Field) *Logger {
	if len(field) == 0 {
		return log
	}

	log.l = log.l.With(field...)
	return log
}

func (log *Logger) Debug(msg string, field ...zap.Field) {
	log.l.Debug(msg, field...)
}

func (log *Logger) Info(msg string, field ...zap.Field) {
	log.l.Info(msg, field...)
}

func (log *Logger) Warn(msg string, field ...zap.Field) {
	log.l.Warn(msg, field...)
}

func (log *Logger) Error(msg string, field ...zap.Field) {
	log.l.Error(msg, field...)
}

func (log *Logger) DPanic(msg string, field ...zap.Field) {
	log.l.DPanic(msg, field...)
}

func (log *Logger) Panic(msg string, field ...zap.Field) {
	log.l.Panic(msg, field...)
}

func (log *Logger) Fatal(msg string, field ...zap.Field) {
	log.l.Fatal(msg, field...)
}
