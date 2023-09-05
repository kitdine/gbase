package log

import (
	"go.uber.org/zap"
)

type Field = zap.Field

var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any
)

func Named(s string) *Logger {
	l := &Logger{
		l: GetLogger().l,
	}
	return l.Named(s)
}

func With(fields ...Field) *Logger {
	l := &Logger{
		l: GetLogger().l,
	}
	return l.With(fields...)
}

func Debug(msg string, fields ...Field) {
	GetLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	GetLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	GetLogger().Error(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	GetLogger().DPanic(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	GetLogger().Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	GetLogger().Fatal(msg, fields...)
}

func (log *Logger) Named(s string) *Logger {
	log.l = log.l.Named(s)
	return log
}

func (log *Logger) With(fields ...Field) *Logger {
	if len(fields) == 0 {
		return log
	}

	log.l = log.l.With(fields...)
	return log
}

func (log *Logger) Debug(msg string, fields ...Field) {
	log.l.Debug(msg, fields...)
}

func (log *Logger) Info(msg string, fields ...Field) {
	log.l.Info(msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...Field) {
	log.l.Warn(msg, fields...)
}

func (log *Logger) Error(msg string, fields ...Field) {
	log.l.Error(msg, fields...)
}

func (log *Logger) DPanic(msg string, fields ...Field) {
	log.l.DPanic(msg, fields...)
}

func (log *Logger) Panic(msg string, fields ...Field) {
	log.l.Panic(msg, fields...)
}

func (log *Logger) Fatal(msg string, fields ...Field) {
	log.l.Fatal(msg, fields...)
}
