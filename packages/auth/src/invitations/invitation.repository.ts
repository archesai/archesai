import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { InvitationEntity } from '@archesai/schemas'

import { createBaseRepository, InvitationTable } from '@archesai/database'
import { InvitationEntitySchema } from '@archesai/schemas'

export const createInvitationRepository = (
  databaseService: DatabaseService
): BaseRepository<
  InvitationEntity,
  (typeof InvitationTable)['$inferInsert'],
  (typeof InvitationTable)['$inferSelect']
> => {
  return createBaseRepository<InvitationEntity>(
    databaseService,
    InvitationTable,
    InvitationEntitySchema
  )
}

export type InvitationRepository = ReturnType<typeof createInvitationRepository>
