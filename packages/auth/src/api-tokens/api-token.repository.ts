import type { DatabaseService } from '@archesai/core'
import type {
  ApiTokenInsertModel,
  ApiTokenSelectModel
} from '@archesai/database'
import type { ApiTokenEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { ApiTokenTable } from '@archesai/database'
import { ApiTokenEntitySchema } from '@archesai/schemas'

export class ApiTokenRepository extends BaseRepository<
  ApiTokenEntity,
  ApiTokenInsertModel,
  ApiTokenSelectModel
> {
  constructor(
    databaseService: DatabaseService<ApiTokenInsertModel, ApiTokenSelectModel>
  ) {
    super(databaseService, ApiTokenTable, ApiTokenEntitySchema)
  }
}
