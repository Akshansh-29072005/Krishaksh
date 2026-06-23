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

	// GCP Cloud Tasks Config
	GCPProjectID        string
	GCPLocation         string
	CloudTasksQueue     string
	CloudTasksWorkerURL string // URL for Cloud Run Jobs/Cloud Functions worker endpoint
	GCSBucket           string // Name of the Google Cloud Storage bucket
}

func LoadConfig() *Config {
	// Load .env only if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it. Using system environment variables.")
	}

	// Cloud Run uses PORT, so prioritize that
	port := getEnvOrDefault("PORT", "")
	if port == "" {
		port = getEnvOrDefault("SERVER_PORT", "8080")
	}

	return &Config{
		ServerPort:          port,
		ServerEnv:           getEnvOrDefault("SERVER_ENV", "development"),
		DatabaseURL:         getEnvOrDefault("DATABASE_URL", ""),
		RedisAddr:           getEnvOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Env:                 getEnvOrDefault("SERVER_ENV", "development"),
		DBMaxConns:          getEnvOrDefault("DB_MAX_CONNS", "20"),
		DBMinConns:          getEnvOrDefault("DB_MIN_CONNS", "2"),
		DBConnMaxLifetime:   getEnvOrDefault("DB_CONN_MAX_LIFETIME", "1h"),
		DBConnMaxIdleTime:   getEnvOrDefault("DB_CONN_MAX_IDLE_TIME", "30m"),
		GCPProjectID:        getEnvOrDefault("GCP_PROJECT_ID", ""),
		GCPLocation:         getEnvOrDefault("GCP_LOCATION", "us-central1"),
		CloudTasksQueue:     getEnvOrDefault("CLOUD_TASKS_QUEUE", "krisho-queue"),
		CloudTasksWorkerURL: getEnvOrDefault("CLOUD_TASKS_WORKER_URL", ""),
		GCSBucket:           getEnvOrDefault("GCS_BUCKET", ""),
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
		// Log specific error to help debugging
		log.Println("error: DATABASE_URL environment variable is empty")
		return errors.New("DATABASE_URL is required for Supabase connection")
	}
	if c.ServerPort == "" {
		return errors.New("server port missing")
	}
	// Redis is optional now (only required if QUEUE_TYPE != cloudtasks)
	if os.Getenv("QUEUE_TYPE") != "cloudtasks" && c.RedisAddr == "" {
		return errors.New("redis addr missing (required when QUEUE_TYPE != cloudtasks)")
	}
	return nil
}
