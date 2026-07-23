package config

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

type AppConfig struct {
	Name string
    Env  string
}

type DatabaseConfig struct {}

type HttpConfig struct {
	Host string
	Port int
}

type SessionConfig struct {
	CookieName string
    TTL         time.Duration
    Secure      bool
    HTTPOnly    bool
}

type RedisConfig struct {}
type AIConfig struct {}
type PaymentConfig struct {}
type LoggerConfig struct {
	LogLevel       Level
}


type Config struct {
	App AppConfig
	Database DatabaseConfig
	Redis RedisConfig
	Http HttpConfig
	Session SessionConfig
	AI AIConfig
	Payment PaymentConfig
	Logger LoggerConfig
}