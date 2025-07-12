import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { LabelTable } from '@archesai/database'
import { LabelEntity } from '@archesai/schemas'

/**
 * Repository for labels.
 */
export class LabelRepository extends BaseRepository<LabelEntity> {
  constructor(databaseService: DatabaseService<LabelEntity>) {
    super(databaseService, LabelTable, LabelEntity)
  }
}
