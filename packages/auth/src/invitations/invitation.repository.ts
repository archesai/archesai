import type { DatabaseService } from '@archesai/core'
import type {
  InvitationInsertModel,
  InvitationSelectModel
} from '@archesai/database'
import type { InvitationEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { InvitationTable } from '@archesai/database'
import { InvitationEntitySchema } from '@archesai/schemas'

export const createInvitationRepository = (
  databaseService: DatabaseService<InvitationInsertModel, InvitationSelectModel>
) => {
  return createBaseRepository<
    InvitationEntity,
    InvitationInsertModel,
    InvitationSelectModel
  >(databaseService, InvitationTable, InvitationEntitySchema)
}

export type InvitationRepository = ReturnType<typeof createInvitationRepository>
