import type { DatabaseService } from '@archesai/database'
import type { ArtifactEntity } from '@archesai/schemas'

import { ArtifactTable, createBaseRepository } from '@archesai/database'
import { ArtifactEntitySchema } from '@archesai/schemas'

export const createArtifactRepository = (databaseService: DatabaseService) => {
  return createBaseRepository<ArtifactEntity>(
    databaseService,
    ArtifactTable,
    ArtifactEntitySchema
  )
}

export type ArtifactRepository = ReturnType<typeof createArtifactRepository>
