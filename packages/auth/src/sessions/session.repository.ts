import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { SessionEntity } from '@archesai/schemas'

import { createBaseRepository, SessionTable } from '@archesai/database'
import { SessionEntitySchema } from '@archesai/schemas'

export const createSessionRepository = (
  databaseService: DatabaseService
): BaseRepository<
  SessionEntity,
  (typeof SessionTable)['$inferInsert'],
  (typeof SessionTable)['$inferSelect']
> => {
  return createBaseRepository<SessionEntity>(
    databaseService,
    SessionTable,
    SessionEntitySchema
  )
}

export type SessionRepository = ReturnType<typeof createSessionRepository>
