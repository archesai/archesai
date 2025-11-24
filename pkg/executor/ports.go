package executor

import (
	"github.com/archesai/archesai/pkg/database"
)

// executorRepository handles Executor persistence
type executorRepository interface {
	database.Repository[Executor]
}
