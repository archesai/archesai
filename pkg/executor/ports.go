package executor

import (
	"github.com/archesai/archesai/apps/studio/generated/core/models"
	"github.com/archesai/archesai/pkg/database"
)

// executorRepository handles Executor persistence
type executorRepository interface {
	database.CRUDRepository[models.Executor] // or database.CRUDRepository[models.Session, uuid.UUID]
}
