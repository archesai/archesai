import type { DatabaseService } from '@archesai/core'
import type {
  VerificationTokenInsertModel,
  VerificationTokenSelectModel
} from '@archesai/database'
import type { VerificationTokenEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { VerificationTokenTable } from '@archesai/database'
import { VerificationTokenEntitySchema } from '@archesai/schemas'

export const createVerificationRepository = (
  databaseService: DatabaseService<
    VerificationTokenInsertModel,
    VerificationTokenSelectModel
  >
) => {
  return createBaseRepository<
    VerificationTokenEntity,
    VerificationTokenInsertModel,
    VerificationTokenSelectModel
  >(databaseService, VerificationTokenTable, VerificationTokenEntitySchema)
}

export type VerificationRepository = ReturnType<
  typeof createVerificationRepository
>
