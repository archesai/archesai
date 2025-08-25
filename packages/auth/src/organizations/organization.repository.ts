import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { OrganizationEntity } from '@archesai/schemas'

import { createBaseRepository, OrganizationTable } from '@archesai/database'
import { OrganizationEntitySchema } from '@archesai/schemas'

export const createOrganizationRepository = (
  databaseService: DatabaseService
): BaseRepository<
  OrganizationEntity,
  (typeof OrganizationTable)['$inferInsert'],
  (typeof OrganizationTable)['$inferSelect']
> => {
  return createBaseRepository<OrganizationEntity>(
    databaseService,
    OrganizationTable,
    OrganizationEntitySchema
  )
}

export type OrganizationRepository = ReturnType<
  typeof createOrganizationRepository
>
