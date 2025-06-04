import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { ApiTokenTable } from '@archesai/database'
import { ApiTokenEntity } from '@archesai/domain'

export class ApiTokenRepository extends BaseRepository<ApiTokenEntity> {
  constructor(databaseService: DatabaseService<ApiTokenEntity>) {
    super(databaseService, ApiTokenTable, ApiTokenEntity)
  }
}
