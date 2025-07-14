import type {
  DrizzleDatabaseService,
  PipelineSelectModel
} from '@archesai/database'
import type { PipelineEntity } from '@archesai/schemas'

import { createBaseRepository, PipelineTable } from '@archesai/database'
import { PipelineEntitySchema } from '@archesai/schemas'

export const createPipelineRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<PipelineEntity, PipelineSelectModel>(
    databaseService,
    PipelineTable,
    PipelineEntitySchema
  )
}

export type PipelineRepository = ReturnType<typeof createPipelineRepository>
