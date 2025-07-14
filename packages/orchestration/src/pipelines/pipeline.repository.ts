import type { DatabaseService } from '@archesai/core'
import type {
  PipelineInsertModel,
  PipelineSelectModel
} from '@archesai/database'
import type { PipelineEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { PipelineTable } from '@archesai/database'
import { PipelineEntitySchema } from '@archesai/schemas'

export const createPipelineRepository = (
  databaseService: DatabaseService<PipelineInsertModel, PipelineSelectModel>
) => {
  return createBaseRepository<
    PipelineEntity,
    PipelineInsertModel,
    PipelineSelectModel
  >(databaseService, PipelineTable, PipelineEntitySchema)
}

export type PipelineRepository = ReturnType<typeof createPipelineRepository>
