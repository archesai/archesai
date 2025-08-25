import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { ArtifactEntity } from '@archesai/schemas'

import { ArtifactTable, createBaseRepository } from '@archesai/database'
import { ArtifactEntitySchema } from '@archesai/schemas'

export const createArtifactRepository = (
  databaseService: DatabaseService
): BaseRepository<
  ArtifactEntity,
  (typeof ArtifactTable)['$inferInsert'],
  (typeof ArtifactTable)['$inferSelect']
> => {
  return createBaseRepository<ArtifactEntity>(
    databaseService,
    ArtifactTable,
    ArtifactEntitySchema
  )
}

export type ArtifactRepository = ReturnType<typeof createArtifactRepository>
