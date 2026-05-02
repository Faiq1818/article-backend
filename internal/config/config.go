package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Database pool
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	// HTTP server
	HTTPReadTimeout     time.Duration
	HTTPWriteTimeout    time.Duration
	HTTPIdleTimeout     time.Duration
	HTTPShutdownTimeout time.Duration

	// JWT
	JWTExpirationHours int

	// Upload
	MaxUploadSizeBytes int64
}

func Load() *Config {
	return &Config{
		DBMaxOpenConns:      getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:      getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime:   getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		DBConnMaxIdleTime:   getEnvDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
		HTTPReadTimeout:     getEnvDuration("HTTP_READ_TIMEOUT", 10*time.Second),
		HTTPWriteTimeout:    getEnvDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		HTTPIdleTimeout:     getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		HTTPShutdownTimeout: getEnvDuration("HTTP_SHUTDOWN_TIMEOUT", 30*time.Second),
		JWTExpirationHours:  getEnvInt("JWT_EXPIRATION_HOURS", 24),
		MaxUploadSizeBytes:  getEnvInt64("MAX_UPLOAD_SIZE_BYTES", 20<<20),
	}
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return v
}

func getEnvInt64(key string, fallback int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	v, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return v
}
