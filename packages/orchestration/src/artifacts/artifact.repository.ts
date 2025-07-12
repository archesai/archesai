import type { DatabaseService } from '@archesai/core'
import type {
  ArtifactInsertModel,
  ArtifactSelectModel
} from '@archesai/database'
import type { ArtifactEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { ArtifactTable } from '@archesai/database'
import { ArtifactEntitySchema } from '@archesai/schemas'

/**
 * Repository for content.
 */
export class ArtifactRepository extends BaseRepository<
  ArtifactEntity,
  ArtifactInsertModel,
  ArtifactSelectModel
> {
  constructor(
    databaseService: DatabaseService<ArtifactInsertModel, ArtifactSelectModel>
  ) {
    super(databaseService, ArtifactTable, ArtifactEntitySchema)
  }
}
