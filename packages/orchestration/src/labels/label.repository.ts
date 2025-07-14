import type { DatabaseService } from '@archesai/core'
import type { LabelInsertModel, LabelSelectModel } from '@archesai/database'
import type { LabelEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { LabelTable } from '@archesai/database'
import { LabelEntitySchema } from '@archesai/schemas'

export const createLabelRepository = (
  databaseService: DatabaseService<LabelInsertModel, LabelSelectModel>
) => {
  return createBaseRepository<LabelEntity, LabelInsertModel, LabelSelectModel>(
    databaseService,
    LabelTable,
    LabelEntitySchema
  )
}

export type LabelRepository = ReturnType<typeof createLabelRepository>
