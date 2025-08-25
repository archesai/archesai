import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { PipelineEntity } from '@archesai/schemas'

import { createBaseRepository, PipelineTable } from '@archesai/database'
import { PipelineEntitySchema } from '@archesai/schemas'

export const createPipelineRepository = (
  databaseService: DatabaseService
): BaseRepository<
  PipelineEntity,
  (typeof PipelineTable)['$inferInsert'],
  (typeof PipelineTable)['$inferSelect']
> => {
  return createBaseRepository<PipelineEntity>(
    databaseService,
    PipelineTable,
    PipelineEntitySchema
  )
}

export type PipelineRepository = ReturnType<typeof createPipelineRepository>
