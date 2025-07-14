import type { DatabaseService } from '@archesai/core'
import type { SessionInsertModel, SessionSelectModel } from '@archesai/database'
import type { SessionEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { SessionTable } from '@archesai/database'
import { SessionEntitySchema } from '@archesai/schemas'

export const createSessionRepository = (
  databaseService: DatabaseService<SessionInsertModel, SessionSelectModel>
) => {
  return createBaseRepository<
    SessionEntity,
    SessionInsertModel,
    SessionSelectModel
  >(databaseService, SessionTable, SessionEntitySchema)
}

export type SessionRepository = ReturnType<typeof createSessionRepository>
