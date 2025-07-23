import type { DatabaseService } from '@archesai/database'
import type { ToolEntity } from '@archesai/schemas'

import { createBaseRepository, ToolTable } from '@archesai/database'
import { ToolEntitySchema } from '@archesai/schemas'

export const createToolRepository = (databaseService: DatabaseService) => {
  return createBaseRepository<ToolEntity>(
    databaseService,
    ToolTable,
    ToolEntitySchema
  )
}

export type ToolRepository = ReturnType<typeof createToolRepository>
