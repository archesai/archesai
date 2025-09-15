package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config wraps the generated ArchesConfig for easy access
type Config struct {
	*ArchesConfig
	v *viper.Viper
}

// Load reads configuration from environment variables and returns a Config.
func Load() (*Config, error) {
	// Start with defaults from generated code
	config := GetDefaultConfig()

	// Setup Viper for environment and file overrides
	v := viper.New()
	setupViper(v)

	// Read config file if exists
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Apply overrides from environment/config to the default config
	applyViperOverrides(config, v)

	return &Config{
		ArchesConfig: config,
		v:            v,
	}, nil
}

// setupViper configures viper for reading config
func setupViper(v *viper.Viper) {
	v.SetConfigName(DefaultConfigName)
	v.SetConfigType(DefaultConfigType)
	for _, path := range ConfigPaths {
		v.AddConfigPath(path)
	}

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}

// applyViperOverrides applies configuration overrides from Viper to the config struct
func applyViperOverrides(config *ArchesConfig, v *viper.Viper) {
	applyAPIOverrides(config, v)
	applyAuthOverrides(config, v)
	applyDatabaseOverrides(config, v)
	applyLoggingOverrides(config, v)
	applyRedisOverrides(config, v)
	applyStorageOverrides(config, v)
	applyInfrastructureOverrides(config, v)
	applyPlatformOverrides(config, v)
	applyIngressOverrides(config, v)
	applyIntelligenceOverrides(config, v)
	applyMonitoringOverrides(config, v)
	applyBillingOverrides(config, v)
}

func applyAPIOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("api.host") {
		config.API.Host = v.GetString("api.host")
	}
	if v.IsSet("api.port") {
		config.API.Port = float32(v.GetInt("api.port"))
	}
	if v.IsSet("api.docs") {
		config.API.Docs = v.GetBool("api.docs")
	}
	if v.IsSet("api.validate") {
		config.API.Validate = v.GetBool("api.validate")
	}
	if v.IsSet("api.environment") {
		config.API.Environment = APIConfigEnvironment(v.GetString("api.environment"))
	}
	if v.IsSet("api.cors.origins") {
		config.API.Cors.Origins = v.GetString("api.cors.origins")
	}
}

func applyAuthOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("auth.enabled") {
		config.Auth.Enabled = v.GetBool("auth.enabled")
	}
	if v.IsSet("auth.local.enabled") {
		config.Auth.Local.Enabled = v.GetBool("auth.local.enabled")
	}
	if v.IsSet("auth.local.jwt_secret") {
		config.Auth.Local.JwtSecret = v.GetString("auth.local.jwt_secret")
	}
	if v.IsSet("auth.local.access_token_ttl") {
		config.Auth.Local.AccessTokenTTL = v.GetString("auth.local.access_token_ttl")
	}
	if v.IsSet("auth.local.refresh_token_ttl") {
		config.Auth.Local.RefreshTokenTTL = v.GetString("auth.local.refresh_token_ttl")
	}
}

func applyDatabaseOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("database.enabled") {
		config.Database.Enabled = v.GetBool("database.enabled")
	}
	if v.IsSet("database.url") {
		config.Database.URL = v.GetString("database.url")
	}
	if v.IsSet("database.type") {
		config.Database.Type = DatabaseConfigType(v.GetString("database.type"))
	}
	if v.IsSet("database.max_conns") {
		config.Database.MaxConns = v.GetInt("database.max_conns")
	}
	if v.IsSet("database.min_conns") {
		config.Database.MinConns = v.GetInt("database.min_conns")
	}
	if v.IsSet("database.conn_max_lifetime") {
		config.Database.ConnMaxLifetime = v.GetString("database.conn_max_lifetime")
	}
	if v.IsSet("database.conn_max_idle_time") {
		config.Database.ConnMaxIdleTime = v.GetString("database.conn_max_idle_time")
	}
	if v.IsSet("database.health_check_period") {
		config.Database.HealthCheckPeriod = v.GetString("database.health_check_period")
	}
	if v.IsSet("database.run_migrations") {
		config.Database.RunMigrations = v.GetBool("database.run_migrations")
	}
}

func applyLoggingOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("logging.level") {
		config.Logging.Level = LoggingConfigLevel(v.GetString("logging.level"))
	}
	if v.IsSet("logging.pretty") {
		config.Logging.Pretty = v.GetBool("logging.pretty")
	}
}

func applyRedisOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("redis.enabled") {
		config.Redis.Enabled = v.GetBool("redis.enabled")
	}
	if v.IsSet("redis.host") {
		config.Redis.Host = v.GetString("redis.host")
	}
	if v.IsSet("redis.port") {
		config.Redis.Port = float32(v.GetInt("redis.port"))
	}
	if v.IsSet("redis.auth") {
		config.Redis.Auth = v.GetString("redis.auth")
	}
}

func applyStorageOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("storage.enabled") {
		config.Storage.Enabled = v.GetBool("storage.enabled")
	}
	if v.IsSet("storage.endpoint") {
		config.Storage.Endpoint = v.GetString("storage.endpoint")
	}
	if v.IsSet("storage.bucket") {
		config.Storage.Bucket = v.GetString("storage.bucket")
	}
	if v.IsSet("storage.access_key") {
		config.Storage.Accesskey = v.GetString("storage.access_key")
	}
	if v.IsSet("storage.secret_key") {
		config.Storage.Secretkey = v.GetString("storage.secret_key")
	}
}

func applyInfrastructureOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("infrastructure.namespace") {
		config.Infrastructure.Namespace = v.GetString("infrastructure.namespace")
	}
	if v.IsSet("infrastructure.development.host_ip") {
		config.Infrastructure.Development.HostIP = v.GetString("infrastructure.development.host_ip")
	}
}

func applyPlatformOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("platform.enabled") {
		config.Platform.Enabled = v.GetBool("platform.enabled")
	}
}

func applyIngressOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("ingress.enabled") {
		config.Ingress.Enabled = v.GetBool("ingress.enabled")
	}
	if v.IsSet("ingress.domain") {
		config.Ingress.Domain = v.GetString("ingress.domain")
	}
}

func applyIntelligenceOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("intelligence.llm.type") {
		config.Intelligence.Llm.Type = LLMConfigType(v.GetString("intelligence.llm.type"))
	}
	if v.IsSet("intelligence.llm.endpoint") {
		config.Intelligence.Llm.Endpoint = v.GetString("intelligence.llm.endpoint")
	}
	if v.IsSet("intelligence.llm.token") {
		config.Intelligence.Llm.Token = v.GetString("intelligence.llm.token")
	}
	if v.IsSet("intelligence.embedding.type") {
		config.Intelligence.Embedding.Type = EmbeddingConfigType(v.GetString("intelligence.embedding.type"))
	}
}

func applyMonitoringOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("monitoring.grafana.enabled") {
		config.Monitoring.Grafana.Enabled = v.GetBool("monitoring.grafana.enabled")
	}
	if v.IsSet("monitoring.loki.enabled") {
		config.Monitoring.Loki.Enabled = v.GetBool("monitoring.loki.enabled")
	}
}

func applyBillingOverrides(config *ArchesConfig, v *viper.Viper) {
	if v.IsSet("billing.enabled") {
		config.Billing.Enabled = v.GetBool("billing.enabled")
	}
	if v.IsSet("billing.stripe.token") {
		config.Billing.Stripe.Token = v.GetString("billing.stripe.token")
	}
	if v.IsSet("billing.stripe.whsec") {
		config.Billing.Stripe.Whsec = v.GetString("billing.stripe.whsec")
	}
}

// Helper methods for backward compatibility and convenience

// GetAllowedOrigins returns CORS allowed origins as a slice
func (c *Config) GetAllowedOrigins() []string {
	origins := c.API.Cors.Origins
	if origins == "" {
		return []string{}
	}
	return strings.Split(origins, ",")
}

// GetServerAddress returns the formatted server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.API.Host, int(c.API.Port))
}

// ParseDuration parses a duration string and returns time.Duration
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// GetJWTSecret returns the JWT secret for authentication
func (c *Config) GetJWTSecret() string {
	return c.Auth.Local.JwtSecret
}

// GetAccessTokenTTL returns the access token TTL as a Duration
func (c *Config) GetAccessTokenTTL() (time.Duration, error) {
	if c.Auth.Local.AccessTokenTTL != "" {
		return ParseDuration(c.Auth.Local.AccessTokenTTL)
	}
	return 15 * time.Minute, nil // default
}

// GetRefreshTokenTTL returns the refresh token TTL as a Duration
func (c *Config) GetRefreshTokenTTL() (time.Duration, error) {
	if c.Auth.Local.RefreshTokenTTL != "" {
		return ParseDuration(c.Auth.Local.RefreshTokenTTL)
	}
	return 168 * time.Hour, nil // default 7 days
}
