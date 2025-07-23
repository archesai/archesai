import type { DatabaseService } from '@archesai/database'
import type { InvitationEntity } from '@archesai/schemas'

import { createBaseRepository, InvitationTable } from '@archesai/database'
import { InvitationEntitySchema } from '@archesai/schemas'

export const createInvitationRepository = (
  databaseService: DatabaseService
) => {
  return createBaseRepository<InvitationEntity>(
    databaseService,
    InvitationTable,
    InvitationEntitySchema
  )
}

export type InvitationRepository = ReturnType<typeof createInvitationRepository>
