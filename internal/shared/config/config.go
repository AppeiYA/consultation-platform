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

type DatabaseConfig struct {
	Host string
	Port int
	User string
	Password string
	Name string
	SSLMode string
	Schema string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxLifetimeMinutes int
}

type HttpConfig struct {
	Host string
	Port int
}

type SessionConfig struct {
	CookieName string
    TTL         time.Duration
    Secure      bool
    HTTPOnly    bool
	SameSite string
}

type RedisConfig struct {
	Address string
	Username string
	Password string
	DB int
}
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