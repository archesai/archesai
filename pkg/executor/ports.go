package executor

import (
	"github.com/archesai/archesai/pkg/database"
	"github.com/archesai/archesai/pkg/executor/schemas"
)

// executorRepository handles Executor persistence
type executorRepository interface {
	database.Repository[schemas.Executor]
}
