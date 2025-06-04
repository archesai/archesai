import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { OrganizationTable } from '@archesai/database'
import { OrganizationEntity } from '@archesai/domain'

/**
 * Repository for handling organizations.
 */
export class OrganizationRepository extends BaseRepository<OrganizationEntity> {
  constructor(databaseService: DatabaseService<OrganizationEntity>) {
    super(databaseService, OrganizationTable, OrganizationEntity)
  }
}
