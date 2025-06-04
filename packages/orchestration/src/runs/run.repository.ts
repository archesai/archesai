import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { RUN_ENTITY_KEY, RunEntity } from '@archesai/domain'

/**
 * Repository for runs.
 */
export class RunRepository extends BaseRepository<RunEntity> {
  constructor(databaseService: DatabaseService<RunEntity>) {
    super(databaseService, RUN_ENTITY_KEY, RunEntity)
  }
}
