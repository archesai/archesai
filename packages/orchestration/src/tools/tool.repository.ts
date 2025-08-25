import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { ToolEntity } from '@archesai/schemas'

import { createBaseRepository, ToolTable } from '@archesai/database'
import { ToolEntitySchema } from '@archesai/schemas'

export const createToolRepository = (
  databaseService: DatabaseService
): BaseRepository<
  ToolEntity,
  (typeof ToolTable)['$inferInsert'],
  (typeof ToolTable)['$inferSelect']
> => {
  return createBaseRepository<ToolEntity>(
    databaseService,
    ToolTable,
    ToolEntitySchema
  )
}

export type ToolRepository = ReturnType<typeof createToolRepository>
