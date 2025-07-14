import type {
  ArtifactSelectModel,
  DrizzleDatabaseService
} from '@archesai/database'
import type { ArtifactEntity } from '@archesai/schemas'

import { ArtifactTable, createBaseRepository } from '@archesai/database'
import { ArtifactEntitySchema } from '@archesai/schemas'

export const createArtifactRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<ArtifactEntity, ArtifactSelectModel>(
    databaseService,
    ArtifactTable,
    ArtifactEntitySchema
  )
}

export type ArtifactRepository = ReturnType<typeof createArtifactRepository>
