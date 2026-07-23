package logger

import (
	"os"
	"sync"

	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.Logger
}

var (
	instance *Logger
	once     sync.Once
)

type Config struct {
	Development bool
}

// New creates a new Logger with the provided config.
func New(cfg *config.Config) *Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.App.Env == "development" {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	atomicLevel := zap.NewAtomicLevelAt(cfg.Logger.LogLevel)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		atomicLevel,
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{zap: zapLogger}
}

// Init initialises the global singleton logger. Safe to call once at startup.
func Init(cfg *config.Config) {
	once.Do(func() {
		instance = New(cfg)
	})
}

// Default returns the global singleton, initialising it with sensible
// defaults (level from LOG_LEVEL env var, JSON output) if Init has not been called yet.
func Default() *Logger {
	Init(&config.Config{})
	return instance
}

// With returns a child logger with the given fields attached to every entry.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{zap: l.zap.With(fields...)}
}

// Named adds a named sub-scope to the logger.
func (l *Logger) Named(name string) *Logger {
	return &Logger{zap: l.zap.Named(name)}
}

// Sync flushes any buffered log entries. Call on application shutdown.
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) { Default().Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { Default().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Default().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Default().Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Default().Fatal(msg, fields...) }
