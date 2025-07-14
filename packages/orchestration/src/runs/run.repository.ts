import type { DatabaseService } from '@archesai/core'
import type { RunInsertModel, RunSelectModel } from '@archesai/database'
import type { RunEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { RunTable } from '@archesai/database'
import { RunEntitySchema } from '@archesai/schemas'

export const createRunRepository = (
  databaseService: DatabaseService<RunInsertModel, RunSelectModel>
) => {
  return createBaseRepository<RunEntity, RunInsertModel, RunSelectModel>(
    databaseService,
    RunTable,
    RunEntitySchema
  )
}

export type RunRepository = ReturnType<typeof createRunRepository>
