import type { DatabaseService } from '@archesai/core'
import type { ArtifactEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { ArtifactTable } from '@archesai/database'
import { ArtifactEntitySchema } from '@archesai/schemas'

/**
 * Repository for content.
 */
export class ArtifactRepository extends BaseRepository<ArtifactEntity> {
  constructor(databaseService: DatabaseService<ArtifactEntity>) {
    super(databaseService, ArtifactTable, ArtifactEntitySchema)
  }
}
