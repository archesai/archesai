import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import {
  VERIFICATION_TOKEN_ENTITY_KEY,
  VerificationTokenEntity
} from '@archesai/domain'

/**
 * Repository for verification tokens.
 */
export class VerificationTokenRepository extends BaseRepository<VerificationTokenEntity> {
  constructor(databaseService: DatabaseService<VerificationTokenEntity>) {
    super(
      databaseService,
      VERIFICATION_TOKEN_ENTITY_KEY,
      VerificationTokenEntity
    )
  }
}
