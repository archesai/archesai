import type {
  DrizzleDatabaseService,
  OrganizationSelectModel
} from '@archesai/database'
import type { OrganizationEntity } from '@archesai/schemas'

import { createBaseRepository, OrganizationTable } from '@archesai/database'
import { OrganizationEntitySchema } from '@archesai/schemas'

export const createOrganizationRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<OrganizationEntity, OrganizationSelectModel>(
    databaseService,
    OrganizationTable,
    OrganizationEntitySchema
  )
}

export type OrganizationRepository = ReturnType<
  typeof createOrganizationRepository
>
