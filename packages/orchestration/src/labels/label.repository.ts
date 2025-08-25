import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { LabelEntity } from '@archesai/schemas'

import { createBaseRepository, LabelTable } from '@archesai/database'
import { LabelEntitySchema } from '@archesai/schemas'

export const createLabelRepository = (
  databaseService: DatabaseService
): BaseRepository<
  LabelEntity,
  (typeof LabelTable)['$inferInsert'],
  (typeof LabelTable)['$inferSelect']
> => {
  return createBaseRepository<LabelEntity>(
    databaseService,
    LabelTable,
    LabelEntitySchema
  )
}

export type LabelRepository = ReturnType<typeof createLabelRepository>
