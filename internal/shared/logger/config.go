package logger

// Config holds the configuration for the logger.
type Config struct {
	Level  string
	Pretty bool
}

// NewConfig returns a Config with default values.
func NewConfig() Config {
	return Config{
		Level:  "info",
		Pretty: false,
	}
}
