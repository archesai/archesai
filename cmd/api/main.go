package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/archesai/archesai/internal/infrastructure/database"
	"github.com/archesai/archesai/internal/infrastructure/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Port           string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	AllowedOrigins []string
	Environment    string
}

func main() {

	// Load configuration
	config := loadConfig()

	// Initialize application with Wire
	ctx := context.Background()
	db, err := initDatabase(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	logger := initLogger()

	// Create Echo instance
	s := server.NewServer(db, logger, &server.Config{
		Port:           config.Port,
		AllowedOrigins: config.AllowedOrigins,
		JWTSecret:      config.JWTSecret,
	})

	// Start server
	s.Start()
}

func initLogger() *zap.Logger {
	config := zap.NewProductionConfig()

	// Configure based on environment
	if os.Getenv("ENVIRONMENT") == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return logger
}

func loadConfig() *Config {
	return &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		DatabaseURL:    getEnvOrDefault("DATABASE_URL", "postgres://localhost/archesai?sslmode=disable"),
		RedisURL:       getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:      getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production"),
		AllowedOrigins: []string{getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000")},
		Environment:    getEnvOrDefault("ENVIRONMENT", "development"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initDatabase(ctx context.Context, databaseURL string) (*database.DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &database.DB{
		Pool: pool,
	}, nil
}
