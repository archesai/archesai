import type {
  DrizzleDatabaseService,
  VerificationTokenSelectModel
} from '@archesai/database'
import type { VerificationTokenEntity } from '@archesai/schemas'

import {
  createBaseRepository,
  VerificationTokenTable
} from '@archesai/database'
import { VerificationTokenEntitySchema } from '@archesai/schemas'

export const createVerificationRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<
    VerificationTokenEntity,
    VerificationTokenSelectModel
  >(databaseService, VerificationTokenTable, VerificationTokenEntitySchema)
}

export type VerificationRepository = ReturnType<
  typeof createVerificationRepository
>
