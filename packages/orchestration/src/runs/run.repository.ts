import type { DatabaseService } from '@archesai/core'
import type { RunInsertModel, RunSelectModel } from '@archesai/database'
import type { RunEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { RunTable } from '@archesai/database'
import { RunEntitySchema } from '@archesai/schemas'

/**
 * Repository for runs.
 */
export class RunRepository extends BaseRepository<
  RunEntity,
  RunInsertModel,
  RunSelectModel
> {
  constructor(
    databaseService: DatabaseService<RunInsertModel, RunSelectModel>
  ) {
    super(databaseService, RunTable, RunEntitySchema)
  }
}
