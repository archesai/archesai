import type { DatabaseService } from '@archesai/core'
import type { LabelInsertModel, LabelSelectModel } from '@archesai/database'
import type { LabelEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { LabelTable } from '@archesai/database'
import { LabelEntitySchema } from '@archesai/schemas'

/**
 * Repository for labels.
 */
export class LabelRepository extends BaseRepository<
  LabelEntity,
  LabelInsertModel,
  LabelSelectModel
> {
  constructor(
    databaseService: DatabaseService<LabelInsertModel, LabelSelectModel>
  ) {
    super(databaseService, LabelTable, LabelEntitySchema)
  }
}
