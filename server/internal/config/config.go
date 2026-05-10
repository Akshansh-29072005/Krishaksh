package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	ServerEnv   string
	DatabaseURL string // Combined DSN for Supabase
	RedisAddr   string
	Env         string

	// Pooling settings for Supabase Safe Pooling
	DBMaxConns        string
	DBMinConns        string
	DBConnMaxLifetime string
	DBConnMaxIdleTime string
}

func LoadConfig() *Config {
	// Load .env only if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it. Using system environment variables.")
	}

	return &Config{
		ServerPort:        getEnvOrDefault("SERVER_PORT", "8080"),
		ServerEnv:         getEnvOrDefault("SERVER_ENV", "development"),
		DatabaseURL:       getEnvOrDefault("DATABASE_URL", ""),
		RedisAddr:         getEnvOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Env:               getEnvOrDefault("SERVER_ENV", "development"),
		DBMaxConns:        getEnvOrDefault("DB_MAX_CONNS", "20"),
		DBMinConns:        getEnvOrDefault("DB_MIN_CONNS", "2"),
		DBConnMaxLifetime: getEnvOrDefault("DB_CONN_MAX_LIFETIME", "1h"),
		DBConnMaxIdleTime: getEnvOrDefault("DB_CONN_MAX_IDLE_TIME", "30m"),
	}
}

func (c *Config) GetDSN() string {
	// For Supabase, we prefer the full DATABASE_URL which includes sslmode=require
	return c.DatabaseURL
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required for Supabase connection")
	}
	if c.ServerPort == "" {
		return errors.New("server port missing")
	}
	if c.RedisAddr == "" {
		return errors.New("redis addr missing")
	}
	return nil
}
