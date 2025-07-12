import type { DatabaseService } from '@archesai/core'
import type {
  InvitationInsertModel,
  InvitationSelectModel
} from '@archesai/database'
import type { InvitationEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { InvitationTable } from '@archesai/database'
import { InvitationEntitySchema } from '@archesai/schemas'

/**
 * Repository for interacting with the invitation entity.
 */
export class InvitationRepository extends BaseRepository<
  InvitationEntity,
  InvitationInsertModel,
  InvitationSelectModel
> {
  constructor(
    databaseService: DatabaseService<
      InvitationInsertModel,
      InvitationSelectModel
    >
  ) {
    super(databaseService, InvitationTable, InvitationEntitySchema)
  }
}
