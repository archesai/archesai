package health

import "errors"

var (
	// ErrDatabaseUnavailable is returned when database is not available
	ErrDatabaseUnavailable = errors.New("database unavailable")
)
