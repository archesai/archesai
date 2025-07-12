import type { DatabaseService } from '@archesai/core'
import type { OrganizationEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { OrganizationTable } from '@archesai/database'
import { OrganizationEntitySchema } from '@archesai/schemas'

/**
 * Repository for handling organizations.
 */
export class OrganizationRepository extends BaseRepository<OrganizationEntity> {
  constructor(databaseService: DatabaseService<OrganizationEntity>) {
    super(databaseService, OrganizationTable, OrganizationEntitySchema)
  }
}
