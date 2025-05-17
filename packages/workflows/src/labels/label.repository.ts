import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { LABEL_ENTITY_KEY, LabelEntity } from '@archesai/domain'

/**
 * Repository for labels.
 */
export class LabelRepository extends BaseRepository<LabelEntity> {
  constructor(databaseService: DatabaseService<LabelEntity>) {
    super(databaseService, LABEL_ENTITY_KEY, LabelEntity)
  }
}
