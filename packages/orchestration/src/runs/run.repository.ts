import type { DatabaseService } from '@archesai/core'
import type { RunEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { RunTable } from '@archesai/database'
import { RunEntitySchema } from '@archesai/schemas'

/**
 * Repository for runs.
 */
export class RunRepository extends BaseRepository<RunEntity> {
  constructor(databaseService: DatabaseService<RunEntity>) {
    super(databaseService, RunTable, RunEntitySchema)
  }
}
