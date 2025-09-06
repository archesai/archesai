// Code generated manually. TODO: Implement proper defaults generation.
package config

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *ArchesConfig {
	return &ArchesConfig{
		Api: APIConfig{
			Cors: CORSConfig{
				Origins: "*",
			},
			Host: "0.0.0.0",
			Email: EmailConfig{
				Enabled: false,
			},
			Validate:    true,
			Port:        3001,
			Docs:        true,
			Environment: Development,
		},
		Database: DatabaseConfig{
			Type: Postgresql,
			Url:  "postgresql://admin:password@127.0.0.1:5432/archesai-db",
		},
		Logging: LoggingConfig{
			Level:  Info,
			Pretty: true,
		},
	}
}
