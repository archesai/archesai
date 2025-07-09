import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { ArtifactTable } from '@archesai/database'
import { ArtifactEntity } from '@archesai/schemas'

/**
 * Repository for content.
 */
export class ArtifactRepository extends BaseRepository<ArtifactEntity> {
  constructor(databaseService: DatabaseService<ArtifactEntity>) {
    super(databaseService, ArtifactTable, ArtifactEntity)
  }
}
