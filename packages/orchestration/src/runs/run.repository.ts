import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { RunEntity } from '@archesai/schemas'

import { createBaseRepository, RunTable } from '@archesai/database'
import { RunEntitySchema } from '@archesai/schemas'

export const createRunRepository = (
  databaseService: DatabaseService
): BaseRepository<
  RunEntity,
  (typeof RunTable)['$inferInsert'],
  (typeof RunTable)['$inferSelect']
> => {
  return createBaseRepository<RunEntity>(
    databaseService,
    RunTable,
    RunEntitySchema
  )
}

export type RunRepository = ReturnType<typeof createRunRepository>
