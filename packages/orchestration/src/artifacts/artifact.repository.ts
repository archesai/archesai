import type { DatabaseService } from '@archesai/core'
import type {
  ArtifactInsertModel,
  ArtifactSelectModel
} from '@archesai/database'
import type { ArtifactEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { ArtifactTable } from '@archesai/database'
import { ArtifactEntitySchema } from '@archesai/schemas'

export const createArtifactRepository = (
  databaseService: DatabaseService<ArtifactInsertModel, ArtifactSelectModel>
) => {
  return createBaseRepository<
    ArtifactEntity,
    ArtifactInsertModel,
    ArtifactSelectModel
  >(databaseService, ArtifactTable, ArtifactEntitySchema)
}

export type ArtifactRepository = ReturnType<typeof createArtifactRepository>
