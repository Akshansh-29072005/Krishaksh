package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	ServerEnv  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	RedisAddr  string
	Env        string
}

func LoadConfig() *Config {
	// Load .env only if it exists (in production, we might rely entirely on system ENVs)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it. Using system environment variables.")
	}

	return &Config{
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
		ServerEnv:  getEnvOrDefault("SERVER_ENV", "development"),
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("DB_NAME", "krishaksh_db"),
		DBSSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		RedisAddr:  getEnvOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Env:        getEnvOrDefault("SERVER_ENV", "development"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (c *Config) Validate() error {
	if c.DBHost == "" || c.DBPort == "" || c.DBUser == "" || c.DBName == "" {
		return errors.New("database configuration is incomplete")
	}
	if c.ServerPort == "" {
		return errors.New("server port missing")
	}
	if c.RedisAddr == "" {
		return errors.New("redis addr missing")
	}
	return nil
}
