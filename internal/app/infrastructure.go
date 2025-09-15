// Package app provides dependency injection and application container management
package app

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/archesai/archesai/internal/auth"
	authrepository "github.com/archesai/archesai/internal/auth/adapters/repository"
	"github.com/archesai/archesai/internal/cache"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/content"
	contentrepo "github.com/archesai/archesai/internal/content/adapters/repository"
	"github.com/archesai/archesai/internal/database"
	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/database/sqlite"
	"github.com/archesai/archesai/internal/events"
	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/organizations"
	orgrepo "github.com/archesai/archesai/internal/organizations/adapters/repository"
	"github.com/archesai/archesai/internal/redis"
	"github.com/archesai/archesai/internal/users"
	usersrepo "github.com/archesai/archesai/internal/users/adapters/repository"
	"github.com/archesai/archesai/internal/workflows"
	workflowrepo "github.com/archesai/archesai/internal/workflows/adapters/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Infrastructure holds all infrastructure components
type Infrastructure struct {
	Logger         *slog.Logger
	Database       database.Database
	EventPublisher events.Publisher
	AuthCache      auth.Cache
	UsersCache     users.Cache
	// Single Redis client shared across components
	redisClient *redis.Client
}

// Repositories holds all domain repositories
type Repositories struct {
	Auth          auth.Repository
	Users         users.Repository
	Organizations organizations.Repository
	Workflows     workflows.Repository
	Content       content.Repository
}

// NewInfrastructure creates all infrastructure components
func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
	// Initialize logger
	loggerCfg := logger.Config{
		Level:  string(cfg.Logging.Level),
		Pretty: cfg.Logging.Pretty,
	}
	log := logger.New(loggerCfg)
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
			// Log warning but continue with in-memory implementations
			log.Warn("failed to connect to Redis, using in-memory alternatives", "error", err)
			infra.EventPublisher = events.NewNoOpPublisher()
			// Use memory cache for auth
			accountCache := cache.NewMemoryCache[auth.Account]()
			sessionCache := cache.NewMemoryCache[auth.Session]()
			infra.AuthCache = auth.NewCacheAdapter(accountCache, sessionCache)
			// Use NoOp cache for users (will be replaced later)
			infra.UsersCache = users.NewNoOpCache()
		} else {
			// Use Redis for all components
			infra.redisClient = redisClient
			infra.EventPublisher = events.NewRedisPublisher(redisClient.GetRedisClient())
			// Use Redis cache for auth
			accountCache := cache.NewRedisCache[auth.Account](redisClient.GetRedisClient(), "auth:account")
			sessionCache := cache.NewRedisCache[auth.Session](redisClient.GetRedisClient(), "auth:session")
			infra.AuthCache = auth.NewCacheAdapter(accountCache, sessionCache)
			// Use NoOp cache for users (will be replaced later)
			infra.UsersCache = users.NewNoOpCache()
		}
	} else {
		// Use in-memory implementations when Redis is disabled
		infra.EventPublisher = events.NewNoOpPublisher()
		// Use memory cache for auth
		accountCache := cache.NewMemoryCache[auth.Account]()
		sessionCache := cache.NewMemoryCache[auth.Session]()
		infra.AuthCache = auth.NewCacheAdapter(accountCache, sessionCache)
		// Use NoOp cache for users (will be replaced later)
		infra.UsersCache = users.NewNoOpCache()
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
		dbType = database.DetectTypeFromURL(cfg.Database.Url)
	}

	repos := &Repositories{}

	switch dbType {
	case database.TypePostgreSQL:
		pool, ok := db.Underlying().(*pgxpool.Pool)
		if !ok || pool == nil {
			return nil, fmt.Errorf("failed to get PostgreSQL connection pool")
		}

		repos.Auth = authrepository.NewPostgresRepository(pool)
		repos.Users = usersrepo.NewPostgresRepository(pool)
		repos.Organizations = orgrepo.NewPostgresRepository(pool)
		repos.Workflows = workflowrepo.NewPostgresRepository(pool)
		repos.Content = contentrepo.NewPostgresRepository(pool)

	case database.TypeSQLite:
		sqlDB, ok := db.Underlying().(*sql.DB)
		if !ok || sqlDB == nil {
			return nil, fmt.Errorf("failed to get SQLite connection")
		}

		repos.Auth = authrepository.NewSQLiteRepository(sqlDB)
		repos.Users = usersrepo.NewSQLiteRepository(sqlDB)
		repos.Organizations = orgrepo.NewSQLiteRepository(sqlDB)
		repos.Workflows = workflowrepo.NewSQLiteRepository(sqlDB)
		repos.Content = contentrepo.NewSQLiteRepository(sqlDB)

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
		dbType = database.DetectTypeFromURL(cfg.Database.Url)
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
