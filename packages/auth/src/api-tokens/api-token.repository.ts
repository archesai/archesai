import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { API_TOKEN_ENTITY_KEY, ApiTokenEntity } from '@archesai/domain'

export class ApiTokenRepository extends BaseRepository<ApiTokenEntity> {
  constructor(databaseService: DatabaseService<ApiTokenEntity>) {
    super(databaseService, API_TOKEN_ENTITY_KEY, ApiTokenEntity)
  }
}
