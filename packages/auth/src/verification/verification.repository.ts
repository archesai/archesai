import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { VerificationEntity } from '@archesai/schemas'

import { createBaseRepository, VerificationTable } from '@archesai/database'
import { VerificationEntitySchema } from '@archesai/schemas'

export const createVerificationRepository = (
  databaseService: DatabaseService
): BaseRepository<
  VerificationEntity,
  (typeof VerificationTable)['$inferInsert'],
  (typeof VerificationTable)['$inferSelect']
> => {
  return createBaseRepository<VerificationEntity>(
    databaseService,
    VerificationTable,
    VerificationEntitySchema
  )
}

export type VerificationRepository = ReturnType<
  typeof createVerificationRepository
>
