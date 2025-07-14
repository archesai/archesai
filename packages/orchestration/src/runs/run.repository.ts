import type { DrizzleDatabaseService, RunSelectModel } from '@archesai/database'
import type { RunEntity } from '@archesai/schemas'

import { createBaseRepository, RunTable } from '@archesai/database'
import { RunEntitySchema } from '@archesai/schemas'

export const createRunRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<RunEntity, RunSelectModel>(
    databaseService,
    RunTable,
    RunEntitySchema
  )
}

export type RunRepository = ReturnType<typeof createRunRepository>
