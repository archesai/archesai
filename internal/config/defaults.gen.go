// Code generated manually. TODO: Implement proper defaults generation.
package config

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *ArchesConfig {
	return &ArchesConfig{
		Api: APIConfig{
			Port:        8080,
			Docs:        true,
			Environment: Development,
		},
		Database: DatabaseConfig{
			Type: Postgresql,
			Url:  "postgres://localhost/archesai",
		},
		Logging: LoggingConfig{
			Level:  Info,
			Pretty: true,
		},
	}
}
