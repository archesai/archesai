import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { VerificationTokenTable } from '@archesai/database'
import { VerificationTokenEntity } from '@archesai/schemas'

/**
 * Repository for verification tokens.
 */
export class VerificationTokenRepository extends BaseRepository<VerificationTokenEntity> {
  constructor(databaseService: DatabaseService<VerificationTokenEntity>) {
    super(databaseService, VerificationTokenTable, VerificationTokenEntity)
  }
}
