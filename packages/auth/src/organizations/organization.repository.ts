import type { DatabaseService } from '@archesai/core'
import type {
  OrganizationInsertModel,
  OrganizationSelectModel
} from '@archesai/database'
import type { OrganizationEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { OrganizationTable } from '@archesai/database'
import { OrganizationEntitySchema } from '@archesai/schemas'

export const createOrganizationRepository = (
  databaseService: DatabaseService<
    OrganizationInsertModel,
    OrganizationSelectModel
  >
) => {
  return createBaseRepository<
    OrganizationEntity,
    OrganizationInsertModel,
    OrganizationSelectModel
  >(databaseService, OrganizationTable, OrganizationEntitySchema)
}

export type OrganizationRepository = ReturnType<
  typeof createOrganizationRepository
>
