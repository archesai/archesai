// Package app provides dependency injection and application container management
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/cache"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/database"
	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/database/sqlite"
	"github.com/archesai/archesai/internal/events"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/invitations"
	"github.com/archesai/archesai/internal/labels"
	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/members"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/pipelines"
	"github.com/archesai/archesai/internal/redis"
	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/tools"
	"github.com/archesai/archesai/internal/users"
)

// Infrastructure holds all infrastructure components.
type Infrastructure struct {
	Logger         *slog.Logger
	Database       *database.Database // Database wrapper for both PostgreSQL and SQLite
	EventPublisher events.Publisher
	AuthCache      cache.Cache[sessions.Session]
	UsersCache     cache.Cache[users.User]
	// Single Redis client shared across components
	redisClient *redis.Client
}

// Repositories holds all domain repositories.
type Repositories struct {
	Accounts      accounts.Repository
	Sessions      sessions.Repository
	Users         users.Repository
	Organizations organizations.Repository
	Pipelines     pipelines.Repository
	Runs          runs.Repository
	Tools         tools.Repository
	Artifacts     artifacts.Repository
	Labels        labels.Repository
	Members       members.Repository
	Invitations   invitations.Repository
	Health        health.Repository
}

// NewInfrastructure creates all infrastructure components.
func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
	// Initialize logger
	log := logger.New(logger.Config{
		Level:  string(cfg.Logging.Level),
		Pretty: cfg.Logging.Pretty,
	})
	slog.SetDefault(log)

	// Initialize database
	var sqlDB *sql.DB
	var pgxPool *pgxpool.Pool
	var dbType database.Type

	// Determine database type
	if cfg.Database.Type != "" {
		dbType = database.ParseTypeFromString(string(cfg.Database.Type))
	} else {
		dbType = database.DetectTypeFromURL(cfg.Database.URL)
	}

	// Connect to database based on type
	switch dbType {
	case database.TypePostgreSQL:
		// For PostgreSQL, create pgxpool and convert to sql.DB
		poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse database URL: %w", err)
		}

		// Configure connection pool
		if cfg.Database.MaxConns > 0 {
			poolConfig.MaxConns = int32(cfg.Database.MaxConns)
		} else {
			poolConfig.MaxConns = 25
		}

		pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
		pgxPool = pool
		sqlDB = stdlib.OpenDBFromPool(pool)
		log.Info("connected to PostgreSQL database")

	case database.TypeSQLite:
		// For SQLite, directly create sql.DB
		db, err := sql.Open("sqlite3", cfg.Database.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
		}
		sqlDB = db
		log.Info("connected to SQLite database")

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create database wrapper
	db := database.NewDatabase(sqlDB, pgxPool, dbType)

	infra := &Infrastructure{
		Logger:   log,
		Database: db,
	}

	// Initialize Redis-based components if enabled
	if cfg.Redis.Enabled {
		redisConfig := &redis.Config{
			Host:         cfg.Redis.Host,
			Port:         int(cfg.Redis.Port),
			Password:     cfg.Redis.Auth,
			DB:           0,
			EnablePubSub: true,
		}

		redisClient, err := redis.NewClient(redisConfig, log)
		if err != nil {
			log.Warn("failed to connect to Redis, using in-memory alternatives", "error", err)
			infra.EventPublisher = events.NewNoOpPublisher()
			infra.AuthCache = cache.NewMemoryCache[sessions.Session]()
			infra.UsersCache = cache.NewMemoryCache[users.User]()
		} else {
			log.Info("connected to redis", "host", cfg.Redis.Host, "port", cfg.Redis.Port)
			infra.redisClient = redisClient
			infra.EventPublisher = events.NewRedisPublisher(redisClient.GetRedisClient())
			infra.AuthCache = cache.NewRedisCache[sessions.Session](redisClient.GetRedisClient(), "auth:session")
			infra.UsersCache = cache.NewRedisCache[users.User](redisClient.GetRedisClient(), "users")
		}
	} else {
		// Use in-memory alternatives when Redis is disabled
		infra.EventPublisher = events.NewNoOpPublisher()
		infra.AuthCache = cache.NewMemoryCache[sessions.Session]()
		infra.UsersCache = cache.NewMemoryCache[users.User]()
	}

	return infra, nil
}

// NewRepositories creates all domain repositories based on database type.
func NewRepositories(infra *Infrastructure) (*Repositories, error) {
	repos := &Repositories{}

	if infra.Database.IsPostgreSQL() {
		// Use pgxpool for PostgreSQL repositories
		pool := infra.Database.PgxPool()

		// Core repositories
		repos.Accounts = accounts.NewPostgresRepository(pool)
		repos.Sessions = sessions.NewPostgresRepository(pool)
		repos.Users = users.NewPostgresRepository(pool)
		repos.Organizations = organizations.NewPostgresRepository(pool)

		// Pipeline repositories
		repos.Pipelines = pipelines.NewPostgresRepository(pool)
		repos.Runs = runs.NewPostgresRepository(pool)
		repos.Tools = tools.NewPostgresRepository(pool)

		// Content repositories
		repos.Artifacts = artifacts.NewPostgresRepository(pool)
		repos.Labels = labels.NewPostgresRepository(pool)

		// Organization-related repositories
		repos.Members = members.NewPostgresRepository(pool)
		repos.Invitations = invitations.NewPostgresRepository(pool)

		// Health repository uses sql.DB
		repos.Health = health.NewPostgresRepository(infra.Database.SQLDB())

	} else {
		// Use sql.DB for SQLite repositories
		db := infra.Database.SQLDB()

		// Core repositories
		repos.Accounts = accounts.NewSQLiteRepository(db)
		repos.Sessions = sessions.NewSQLiteRepository(db)
		repos.Users = users.NewSQLiteRepository(db)
		repos.Organizations = organizations.NewSQLiteRepository(db)

		// Pipeline repositories
		repos.Pipelines = pipelines.NewSQLiteRepository(db)
		repos.Runs = runs.NewSQLiteRepository(db)
		repos.Tools = tools.NewSQLiteRepository(db)

		// Content repositories
		repos.Artifacts = artifacts.NewSQLiteRepository(db)
		repos.Labels = labels.NewSQLiteRepository(db)

		// Organization-related repositories
		repos.Members = members.NewSQLiteRepository(db)
		repos.Invitations = invitations.NewSQLiteRepository(db)

		// Health repository
		repos.Health = health.NewSQLiteRepository(db)

	}

	return repos, nil
}

// GetQueries returns database-specific query objects.
func GetQueries(
	infra *Infrastructure,
) (pgQueries *postgresql.Queries, sqliteQueries *sqlite.Queries) {
	if infra.Database.IsPostgreSQL() {
		pgQueries = postgresql.New(infra.Database.PgxPool())
	} else {
		sqliteQueries = sqlite.New(infra.Database.SQLDB())
	}

	return pgQueries, sqliteQueries
}

// Close cleans up all infrastructure resources
func (i *Infrastructure) Close() error {
	if i.Database != nil && i.Database.SQLDB() != nil {
		if err := i.Database.SQLDB().Close(); err != nil {
			i.Logger.Error("failed to close database", "error", err)
			return err
		}
	}

	if i.Database != nil && i.Database.PgxPool() != nil {
		i.Database.PgxPool().Close()
	}

	if i.redisClient != nil {
		if err := i.redisClient.Close(); err != nil {
			i.Logger.Error("failed to close redis", "error", err)
			return err
		}
	}

	return nil
}
