package config

import (
	"os"
	"time"
)

type Config struct {
	DB         string
	DBConn     string
	JWTKey     string
	OriginHost string
	ServerHost string
	DBTimeout  time.Duration
}

func New() *Config {
	return &Config{
		DB:         getEnv("DB", "postgres"),
		DBConn:     getEnv("DB_CONN", "postgres://postgres:postgrespw@localhost:32768/gochat?sslmode=disable"),
		JWTKey:     getEnv("JWT_KEY", "secret"),
		OriginHost: getEnv("ORIGIN_HOST", "http://localhost:3000"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0:8080"),
		DBTimeout:  time.Duration(2) * time.Second,
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
