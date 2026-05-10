package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aarcsx/krisho-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func ConnectDB(cfg *config.Config) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dsn := cfg.GetDSN()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing db config: %w", err)
	}

	// Supabase Pooled Connection settings
	maxConns, err := strconv.Atoi(cfg.DBMaxConns)
	if err != nil {
		maxConns = 20 // Fallback
	}

	minConns, err := strconv.Atoi(cfg.DBMinConns)
	if err != nil {
		minConns = 2 // Fallback
	}

	maxLifetime, err := time.ParseDuration(cfg.DBConnMaxLifetime)
	if err != nil {
		maxLifetime = time.Hour
	}

	maxIdleTime, err := time.ParseDuration(cfg.DBConnMaxIdleTime)
	if err != nil {
		maxIdleTime = 30 * time.Minute
	}

	poolConfig.MaxConns = int32(maxConns)
	poolConfig.MinConns = int32(minConns)
	poolConfig.MaxConnLifetime = maxLifetime
	poolConfig.MaxConnIdleTime = maxIdleTime

	// Retry logic for production-safe startup
	var pool *pgxpool.Pool
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				break
			}
			pool.Close()
		}
		log.Printf("Connecting to Supabase (attempt %d/5) failed: %v. Retrying in 5s...", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Supabase after 5 attempts: %w", err)
	}

	log.Println("Successfully connected to Supabase PostgreSQL")

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("PostgreSQL connection pool closed gracefully")
	}
}
