import type {
  DatabaseService,
  VerificationSelectModel
} from '@archesai/database'
import type { VerificationEntity } from '@archesai/schemas'

import { createBaseRepository, VerificationTable } from '@archesai/database'
import { VerificationEntitySchema } from '@archesai/schemas'

export const createVerificationRepository = (
  databaseService: DatabaseService
) => {
  return createBaseRepository<VerificationEntity, VerificationSelectModel>(
    databaseService,
    VerificationTable,
    VerificationEntitySchema
  )
}

export type VerificationRepository = ReturnType<
  typeof createVerificationRepository
>
