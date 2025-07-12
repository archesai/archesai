import type { DatabaseService } from '@archesai/core'
import type {
  OrganizationInsertModel,
  OrganizationSelectModel
} from '@archesai/database'
import type { OrganizationEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { OrganizationTable } from '@archesai/database'
import { OrganizationEntitySchema } from '@archesai/schemas'

/**
 * Repository for handling organizations.
 */
export class OrganizationRepository extends BaseRepository<
  OrganizationEntity,
  OrganizationInsertModel,
  OrganizationSelectModel
> {
  constructor(
    databaseService: DatabaseService<
      OrganizationInsertModel,
      OrganizationSelectModel
    >
  ) {
    super(databaseService, OrganizationTable, OrganizationEntitySchema)
  }
}
