import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { CONTENT_ENTITY_KEY, ContentEntity } from '@archesai/domain'

/**
 * Repository for content.
 */
export class ContentRepository extends BaseRepository<ContentEntity> {
  constructor(databaseService: DatabaseService<ContentEntity>) {
    super(databaseService, CONTENT_ENTITY_KEY, ContentEntity)
  }
}
