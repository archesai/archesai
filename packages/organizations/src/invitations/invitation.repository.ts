import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { INVITATION_ENTITY_KEY, InvitationEntity } from '@archesai/domain'

/**
 * Repository for interacting with the invitation entity.
 */
export class InvitationRepository extends BaseRepository<InvitationEntity> {
  constructor(databaseService: DatabaseService<InvitationEntity>) {
    super(databaseService, INVITATION_ENTITY_KEY, InvitationEntity)
  }
}
