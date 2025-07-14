import type {
  DrizzleDatabaseService,
  InvitationSelectModel
} from '@archesai/database'
import type { InvitationEntity } from '@archesai/schemas'

import { createBaseRepository, InvitationTable } from '@archesai/database'
import { InvitationEntitySchema } from '@archesai/schemas'

export const createInvitationRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<InvitationEntity, InvitationSelectModel>(
    databaseService,
    InvitationTable,
    InvitationEntitySchema
  )
}

export type InvitationRepository = ReturnType<typeof createInvitationRepository>
