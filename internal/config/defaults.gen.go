package config

import (
	"github.com/archesai/archesai/internal/config/generated/api"
)

// GetDefaultConfig returns a new ArchesConfig with all default values from OpenAPI schema
func GetDefaultConfig() *api.ArchesConfig {

	config := &api.ArchesConfig{
		Api: api.APIConfig{
			Cors: api.CORSConfig{
				Origins: "https://platform.archesai.dev",
			},
			Docs: true,
			Email: api.EmailConfig{
				Enabled: false,
			},
			Environment: "development",
			Host:        "0.0.0.0",
			Image: api.ImageConfig{
				PullPolicy: "IfNotPresent",
				Tag:        "latest",
			},
			Port: 3001,
			Resources: api.ResourceConfig{
				Limits:   api.ResourceLimits{},
				Requests: api.ResourceRequests{},
			},
			Validate: true,
		},
		Auth: api.AuthConfig{
			Enabled: true,
			Firebase: api.FirebaseAuth{
				Enabled: false,
			},
			Local: api.LocalAuth{
				AccessTokenTtl:  "15m",
				Enabled:         true,
				JwtSecret:       "change-me-in-production",
				RefreshTokenTtl: "168h",
			},
			Twitter: api.TwitterAuth{
				Enabled: false,
			},
		},
		Billing: api.BillingConfig{
			Enabled: false,
			Stripe:  api.StripeConfig{},
		},
		Database: api.DatabaseConfig{
			Auth: api.DatabaseAuth{
				Database: "archesai-db",
				Password: "password",
			},
			Enabled: true,
			Image: api.ImageConfig{
				PullPolicy: "IfNotPresent",
				Tag:        "latest",
			},
			Managed:  false,
			MaxConns: 25,
			MinConns: 5,
			Persistence: api.PersistenceConfig{
				Enabled: true,
				Size:    "10Gi",
			},
			Resources: api.ResourceConfig{
				Limits:   api.ResourceLimits{},
				Requests: api.ResourceRequests{},
			},
			RunMigrations: false,
			Type:          "postgresql",
			Url:           "postgresql://admin:password@localhost:5432/archesai-db?schema=public",
		},
		Infrastructure: api.InfrastructureConfig{
			Development: api.DevelopmentConfig{
				Api: api.DevServiceConfig{
					Enabled: false,
				},
				HostIP: "172.18.0.1",
				Loki: api.DevServiceConfig{
					Enabled: false,
				},
				Platform: api.DevServiceConfig{
					Enabled: false,
				},
				Postgres: api.DevServiceConfig{
					Enabled: false,
				},
				Redis: api.DevServiceConfig{
					Enabled: false,
				},
			},
			Images: api.ImagesConfig{
				ImagePullSecrets: []string{},
				ImageRegistry:    "",
			},
			Migrations: api.MigrationsConfig{
				Enabled: false,
			},
			Namespace: "arches-system",
			ServiceAccount: api.ServiceAccountConfig{
				Create: true,
				Name:   "",
			},
		},
		Ingress: api.IngressConfig{
			Domain:  "archesai.dev",
			Enabled: false,
			Tls: api.TLSConfig{
				Enabled:    true,
				Issuer:     "letsencrypt-staging",
				SecretName: "archesai-tls",
			},
		},
		Intelligence: api.IntelligenceConfig{
			Embedding: api.EmbeddingConfig{
				Type: "ollama",
			},
			Llm: api.LLMConfig{
				Type: "ollama",
			},
			Runpod: api.RunPodConfig{
				Enabled: false,
			},
			Scraper: api.ScraperConfig{
				Enabled: false,
				Image: api.ImageConfig{
					PullPolicy: "IfNotPresent",
					Tag:        "latest",
				},
				Managed: false,
				Resources: api.ResourceConfig{
					Limits:   api.ResourceLimits{},
					Requests: api.ResourceRequests{},
				},
			},
			Speech: api.SpeechConfig{
				Enabled: false,
			},
			Unstructured: api.UnstructuredConfig{
				Enabled: false,
				Image: api.ImageConfig{
					PullPolicy: "IfNotPresent",
					Tag:        "latest",
				},
				Managed: false,
				Resources: api.ResourceConfig{
					Limits:   api.ResourceLimits{},
					Requests: api.ResourceRequests{},
				},
			},
		},
		Logging: api.LoggingConfig{
			Level:  "info",
			Pretty: false,
		},
		Monitoring: api.MonitoringConfig{
			Grafana: api.GrafanaConfig{
				Enabled: false,
				Image: api.ImageConfig{
					PullPolicy: "IfNotPresent",
					Tag:        "latest",
				},
				Managed: false,
				Resources: api.ResourceConfig{
					Limits:   api.ResourceLimits{},
					Requests: api.ResourceRequests{},
				},
			},
			Loki: api.LokiConfig{
				Enabled: false,
				Host:    "http://localhost:3100",
				Image: api.ImageConfig{
					PullPolicy: "IfNotPresent",
					Tag:        "latest",
				},
				Managed: false,
				Resources: api.ResourceConfig{
					Limits:   api.ResourceLimits{},
					Requests: api.ResourceRequests{},
				},
			},
		},
		Platform: api.PlatformConfig{
			Enabled: false,
			Host:    "localhost",
			Image: api.ImageConfig{
				PullPolicy: "IfNotPresent",
				Tag:        "latest",
			},
			Managed: false,
			Resources: api.ResourceConfig{
				Limits:   api.ResourceLimits{},
				Requests: api.ResourceRequests{},
			},
		},
		Redis: api.RedisConfig{
			Auth:    "password",
			Enabled: false,
			Host:    "localhost",
			Image: api.ImageConfig{
				PullPolicy: "IfNotPresent",
				Tag:        "latest",
			},
			Managed: false,
			Persistence: api.PersistenceConfig{
				Enabled: true,
				Size:    "10Gi",
			},
			Port: 6379,
			Resources: api.ResourceConfig{
				Limits:   api.ResourceLimits{},
				Requests: api.ResourceRequests{},
			},
		},
		Storage: api.StorageConfig{
			Accesskey: "minioadmin",
			Bucket:    "archesai",
			Enabled:   false,
			Endpoint:  "http://localhost:9000",
			Image: api.ImageConfig{
				PullPolicy: "IfNotPresent",
				Tag:        "latest",
			},
			Managed: false,
			Persistence: api.PersistenceConfig{
				Enabled: true,
				Size:    "10Gi",
			},
			Resources: api.ResourceConfig{
				Limits:   api.ResourceLimits{},
				Requests: api.ResourceRequests{},
			},
			Secretkey: "minioadmin",
		},
	}

	return config
}

// GetDefaultConfigWithOverrides returns a config with defaults and applies common overrides
func GetDefaultConfigWithOverrides() *api.ArchesConfig {
	config := GetDefaultConfig()

	// Apply common development overrides
	// These can be customized based on environment

	return config
}
