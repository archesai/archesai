// Package app provides dependency injection and application container management
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/archesai/archesai/internal/core/aggregates"
	ports "github.com/archesai/archesai/internal/core/ports/repositories"
	"github.com/archesai/archesai/internal/core/valueobjects"
	"github.com/archesai/archesai/internal/infrastructure/cache"
	"github.com/archesai/archesai/internal/infrastructure/config"
	"github.com/archesai/archesai/internal/infrastructure/events"
	database "github.com/archesai/archesai/internal/infrastructure/persistence"
	"github.com/archesai/archesai/internal/infrastructure/persistence/postgres/repositories"

	"github.com/archesai/archesai/internal/infrastructure/redis"
	"github.com/archesai/archesai/internal/shared/logger"
)

// Infrastructure holds all infrastructure components.
type Infrastructure struct {
	Logger         *slog.Logger
	Database       *database.Database // Database wrapper for both PostgreSQL and SQLite
	EventPublisher events.Publisher
	AuthCache      cache.Cache[valueobjects.Session]
	UsersCache     cache.Cache[aggregates.User]
	// Single Redis client shared across components
	redisClient *redis.Client
}

// Repositories holds all domain repositories.
type Repositories struct {
	Accounts      ports.AccountRepository
	Sessions      ports.SessionRepository
	Users         ports.UserRepository
	Organizations ports.OrganizationRepository
	Pipelines     ports.PipelineRepository
	Runs          ports.RunRepository
	Tools         ports.ToolRepository
	Artifacts     ports.ArtifactRepository
	Labels        ports.LabelRepository
	Members       ports.MemberRepository
	Invitations   ports.InvitationRepository
	Health        ports.HealthCheckRepository
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
			poolConfig.MaxConns = cfg.Database.MaxConns
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
			infra.AuthCache = cache.NewMemoryCache[valueobjects.Session]()
			infra.UsersCache = cache.NewMemoryCache[aggregates.User]()
		} else {
			log.Info("connected to redis", "host", cfg.Redis.Host, "port", cfg.Redis.Port)
			infra.redisClient = redisClient
			infra.EventPublisher = events.NewRedisPublisher(redisClient.GetRedisClient())
			infra.AuthCache = cache.NewRedisCache[valueobjects.Session](redisClient.GetRedisClient(), "auth:session")
			infra.UsersCache = cache.NewRedisCache[aggregates.User](redisClient.GetRedisClient(), "users")
		}
	} else {
		// Use in-memory alternatives when Redis is disabled
		infra.EventPublisher = events.NewNoOpPublisher()
		infra.AuthCache = cache.NewMemoryCache[valueobjects.Session]()
		infra.UsersCache = cache.NewMemoryCache[aggregates.User]()
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
		repos.Accounts = repositories.NewPostgresAccountRepository(pool)
		repos.Sessions = repositories.NewPostgresSessionRepository(pool)
		repos.Users = repositories.NewPostgresUserRepository(pool)
		repos.Organizations = repositories.NewPostgresOrganizationRepository(pool)

		// Pipeline repositories
		repos.Pipelines = repositories.NewPostgresPipelineRepository(pool)
		repos.Runs = repositories.NewPostgresRunRepository(pool)
		repos.Tools = repositories.NewPostgresToolRepository(pool)

		// Content repositories
		repos.Artifacts = repositories.NewPostgresArtifactRepository(pool)
		repos.Labels = repositories.NewPostgresLabelRepository(pool)

		// Organization-related repositories
		repos.Members = repositories.NewPostgresMemberRepository(pool)
		repos.Invitations = repositories.NewPostgresInvitationRepository(pool)

		// Health repository uses sql.DB
		// repos.Health = repositories.NewPostgresHealthRepository(infra.Database.SQLDB())

	} else {
		// // Use sql.DB for SQLite repositories
		// db := infra.Database.SQLDB()

		// // Core repositories
		// repos.Accounts = sqlite.NewSQLiteAccountsRepository(db)
		// repos.Sessions = auth.NewSQLiteRepository(db)
		// repos.Users = sqlite.NewSQLiteUsersRepository(db)
		// repos.Organizations = sqlite.NewSQLiteOrganizationsRepository(db)

		// // Pipeline repositories
		// repos.Pipelines = sqlite.NewSQLitePipelinesRepository(db)
		// repos.Runs = sqlite.NewSQLiteRunsRepository(db)
		// repos.Tools = sqlite.NewSQLiteToolsRepository(db)

		// // Content repositories
		// repos.Artifacts = sqlite.NewSQLiteArtifactsRepository(db)
		// repos.Labels = sqlite.NewSQLiteLabelsRepository(db)

		// // Organization-related repositories
		// repos.Members = sqlite.NewSQLiteMembersRepository(db)
		// repos.Invitations = sqlite.NewSQLiteInvitationsRepository(db)

		// // Health repository
		// repos.Health = sqlite.NewSQLiteHealthRepository(db)
		return nil, fmt.Errorf("SQLite repositories are not yet implemented")

	}

	return repos, nil
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
