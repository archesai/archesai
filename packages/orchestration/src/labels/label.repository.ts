import type { DatabaseService } from '@archesai/database'
import type { LabelEntity } from '@archesai/schemas'

import { createBaseRepository, LabelTable } from '@archesai/database'
import { LabelEntitySchema } from '@archesai/schemas'

export const createLabelRepository = (databaseService: DatabaseService) => {
  return createBaseRepository<LabelEntity>(
    databaseService,
    LabelTable,
    LabelEntitySchema
  )
}

export type LabelRepository = ReturnType<typeof createLabelRepository>
