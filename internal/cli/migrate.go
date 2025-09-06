package cli

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Load file source driver
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long: `Manage database migrations for the ArchesAI platform.

This command allows you to apply, rollback, and check the status
of database migrations.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE:  runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	RunE:  runMigrateDown,
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE:  runMigrateStatus,
}

var (
	migrationPath string
	databaseURL   string
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)

	// Flags for migration command
	migrateCmd.PersistentFlags().StringVar(&migrationPath, "path", "internal/infrastructure/database/migrations", "Path to migration files")
	migrateCmd.PersistentFlags().StringVar(&databaseURL, "database-url", "", "Database connection URL")

	// Bind to viper
	if err := viper.BindPFlag("migration.path", migrateCmd.PersistentFlags().Lookup("path")); err != nil {
		log.Fatalf("Failed to bind path flag: %v", err)
	}
	if err := viper.BindPFlag("database.url", migrateCmd.PersistentFlags().Lookup("database-url")); err != nil {
		log.Fatalf("Failed to bind database-url flag: %v", err)
	}
}

func getMigrator() (*migrate.Migrate, error) {
	dbURL := viper.GetString("database.url")
	if dbURL == "" {
		return nil, fmt.Errorf("database URL not configured")
	}

	path := viper.GetString("migration.path")

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migrator
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return m, nil
}

func runMigrateUp(_ *cobra.Command, _ []string) error {
	m, err := getMigrator()
	if err != nil {
		return err
	}
	defer func() {
		_, _ = m.Close()
	}()

	log.Println("Running migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}

func runMigrateDown(_ *cobra.Command, _ []string) error {
	m, err := getMigrator()
	if err != nil {
		return err
	}
	defer func() {
		_, _ = m.Close()
	}()

	log.Println("Rolling back last migration...")
	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	log.Println("✅ Migration rolled back successfully")
	return nil
}

func runMigrateStatus(_ *cobra.Command, _ []string) error {
	m, err := getMigrator()
	if err != nil {
		return err
	}
	defer func() {
		_, _ = m.Close()
	}()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations have been applied yet")
	} else {
		log.Printf("Current migration version: %d", version)
		if dirty {
			log.Println("⚠️  Database is in a dirty state")
		}
	}

	return nil
}
