import type {
  DrizzleDatabaseService,
  LabelSelectModel
} from '@archesai/database'
import type { LabelEntity } from '@archesai/schemas'

import { createBaseRepository, LabelTable } from '@archesai/database'
import { LabelEntitySchema } from '@archesai/schemas'

export const createLabelRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<LabelEntity, LabelSelectModel>(
    databaseService,
    LabelTable,
    LabelEntitySchema
  )
}

export type LabelRepository = ReturnType<typeof createLabelRepository>
