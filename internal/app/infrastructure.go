// Package app provides dependency injection and application container management
package app

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/cache"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/database"
	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/database/sqlite"
	"github.com/archesai/archesai/internal/events"
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
	"github.com/jackc/pgx/v5/pgxpool"
)

// Infrastructure holds all infrastructure components
type Infrastructure struct {
	Logger         *slog.Logger
	Database       database.Database
	EventPublisher events.Publisher
	AuthCache      cache.Cache[sessions.Session]
	UsersCache     cache.Cache[users.User]
	// Single Redis client shared across components
	redisClient *redis.Client
}

// Repositories holds all domain repositories
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
}

// NewInfrastructure creates all infrastructure components
func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {

	// Initialize logger
	loggerCfg := logger.Config{
		Level:  string(cfg.Logging.Level),
		Pretty: cfg.Logging.Pretty,
	}
	var log *slog.Logger
	if loggerCfg.Pretty {
		log = logger.NewPretty(loggerCfg)
	} else {
		log = logger.New(loggerCfg)
	}
	slog.SetDefault(log)

	// Initialize database
	dbFactory := database.NewFactory(log)
	db, err := dbFactory.Create(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

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
		log.Info("redis is disabled, using in-memory alternatives")
		infra.EventPublisher = events.NewNoOpPublisher()
		infra.AuthCache = cache.NewNoOpCache[sessions.Session]()
		infra.UsersCache = cache.NewNoOpCache[users.User]()
	}

	return infra, nil
}

// NewRepositories creates all domain repositories based on database type
func NewRepositories(db database.Database, cfg *config.Config) (*Repositories, error) {
	// Determine database type
	var dbType database.Type
	if cfg.Database.Type != "" {
		dbType = database.ParseTypeFromString(string(cfg.Database.Type))
	} else {
		dbType = database.DetectTypeFromURL(cfg.Database.URL)
	}

	repos := &Repositories{}

	switch dbType {
	case database.TypePostgreSQL:
		pool, ok := db.Underlying().(*pgxpool.Pool)
		if !ok || pool == nil {
			return nil, fmt.Errorf("failed to get PostgreSQL connection pool")
		}

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

	case database.TypeSQLite:
		sqlDB, ok := db.Underlying().(*sql.DB)
		if !ok || sqlDB == nil {
			return nil, fmt.Errorf("failed to get SQLite connection")
		}

		// Core repositories
		repos.Accounts = accounts.NewSQLiteRepository(sqlDB)
		repos.Sessions = sessions.NewSQLiteRepository(sqlDB)
		repos.Users = users.NewSQLiteRepository(sqlDB)
		repos.Organizations = organizations.NewSQLiteRepository(sqlDB)

		// Pipeline repositories
		repos.Pipelines = pipelines.NewSQLiteRepository(sqlDB)
		repos.Runs = runs.NewSQLiteRepository(sqlDB)
		repos.Tools = tools.NewSQLiteRepository(sqlDB)

		// Content repositories
		repos.Artifacts = artifacts.NewSQLiteRepository(sqlDB)
		repos.Labels = labels.NewSQLiteRepository(sqlDB)

		// Organization-related repositories
		repos.Members = members.NewSQLiteRepository(sqlDB)
		repos.Invitations = invitations.NewSQLiteRepository(sqlDB)

	default:
		return nil, fmt.Errorf("unsupported database type: %v", dbType)
	}

	return repos, nil
}

// GetQueries returns database-specific query objects
func GetQueries(db database.Database, cfg *config.Config) (pgQueries *postgresql.Queries, sqliteQueries *sqlite.Queries) {
	// Determine database type
	var dbType database.Type
	if cfg.Database.Type != "" {
		dbType = database.ParseTypeFromString(string(cfg.Database.Type))
	} else {
		dbType = database.DetectTypeFromURL(cfg.Database.URL)
	}

	switch dbType {
	case database.TypePostgreSQL:
		if pool, ok := db.Underlying().(*pgxpool.Pool); ok && pool != nil {
			pgQueries = postgresql.New(pool)
		}
	case database.TypeSQLite:
		if sqlDB, ok := db.Underlying().(*sql.DB); ok && sqlDB != nil {
			sqliteQueries = sqlite.New(sqlDB)
		}
	}

	return pgQueries, sqliteQueries
}

// Close cleans up all infrastructure resources
func (i *Infrastructure) Close() error {
	// Close database
	if i.Database != nil {
		if err := i.Database.Close(); err != nil {
			i.Logger.Error("failed to close database", "error", err)
		}
	}

	// Close Redis if it exists
	if i.redisClient != nil {
		if err := i.redisClient.Close(); err != nil {
			i.Logger.Error("failed to close Redis", "error", err)
		}
	}

	return nil
}
