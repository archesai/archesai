import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { RunTable } from '@archesai/database'
import { RunEntity } from '@archesai/domain'

/**
 * Repository for runs.
 */
export class RunRepository extends BaseRepository<RunEntity> {
  constructor(databaseService: DatabaseService<RunEntity>) {
    super(databaseService, RunTable, RunEntity)
  }
}
