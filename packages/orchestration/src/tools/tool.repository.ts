import type {
  DrizzleDatabaseService,
  ToolSelectModel
} from '@archesai/database'
import type { ToolEntity } from '@archesai/schemas'

import { createBaseRepository, ToolTable } from '@archesai/database'
import { ToolEntitySchema } from '@archesai/schemas'

export const createToolRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<ToolEntity, ToolSelectModel>(
    databaseService,
    ToolTable,
    ToolEntitySchema
  )
}

export type ToolRepository = ReturnType<typeof createToolRepository>
