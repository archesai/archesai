import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { ARTIFACT_ENTITY_KEY, ArtifactEntity } from '@archesai/domain'

/**
 * Repository for content.
 */
export class ArtifactRepository extends BaseRepository<ArtifactEntity> {
  constructor(databaseService: DatabaseService<ArtifactEntity>) {
    super(databaseService, ARTIFACT_ENTITY_KEY, ArtifactEntity)
  }
}
