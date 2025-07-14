import type { DatabaseService } from '@archesai/core'
import type { ToolInsertModel, ToolSelectModel } from '@archesai/database'
import type { ToolEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { ToolTable } from '@archesai/database'
import { ToolEntitySchema } from '@archesai/schemas'

export const createToolRepository = (
  databaseService: DatabaseService<ToolInsertModel, ToolSelectModel>
) => {
  return createBaseRepository<ToolEntity, ToolInsertModel, ToolSelectModel>(
    databaseService,
    ToolTable,
    ToolEntitySchema
  )
}

export type ToolRepository = ReturnType<typeof createToolRepository>
