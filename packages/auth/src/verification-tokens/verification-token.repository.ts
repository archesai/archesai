import type { DatabaseService } from '@archesai/core'
import type {
  VerificationTokenInsertModel,
  VerificationTokenSelectModel
} from '@archesai/database'
import type { VerificationTokenEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { VerificationTokenTable } from '@archesai/database'
import { VerificationTokenEntitySchema } from '@archesai/schemas'

/**
 * Repository for verification tokens.
 */
export class VerificationTokenRepository extends BaseRepository<
  VerificationTokenEntity,
  VerificationTokenInsertModel,
  VerificationTokenSelectModel
> {
  constructor(
    databaseService: DatabaseService<
      VerificationTokenEntity,
      VerificationTokenInsertModel,
      VerificationTokenSelectModel
    >
  ) {
    super(
      databaseService,
      VerificationTokenTable,
      VerificationTokenEntitySchema
    )
  }
}
