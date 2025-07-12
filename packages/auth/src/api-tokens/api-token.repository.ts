import type { DatabaseService } from '@archesai/core'
import type { ApiTokenEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { ApiTokenTable } from '@archesai/database'
import { ApiTokenEntitySchema } from '@archesai/schemas'

export class ApiTokenRepository extends BaseRepository<ApiTokenEntity> {
  constructor(databaseService: DatabaseService<ApiTokenEntity>) {
    super(databaseService, ApiTokenTable, ApiTokenEntitySchema)
  }
}
