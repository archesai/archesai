import type { DatabaseService, SessionSelectModel } from '@archesai/database'
import type { SessionEntity } from '@archesai/schemas'

import { createBaseRepository, SessionTable } from '@archesai/database'
import { SessionEntitySchema } from '@archesai/schemas'

export const createSessionRepository = (databaseService: DatabaseService) => {
  return createBaseRepository<SessionEntity, SessionSelectModel>(
    databaseService,
    SessionTable,
    SessionEntitySchema
  )
}

export type SessionRepository = ReturnType<typeof createSessionRepository>
