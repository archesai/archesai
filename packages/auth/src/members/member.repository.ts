import type { DatabaseService } from '@archesai/core'
import type { MemberInsertModel, MemberSelectModel } from '@archesai/database'
import type { MemberEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { MemberTable } from '@archesai/database'
import { MemberEntitySchema } from '@archesai/schemas'

/**
 * Repository for interacting with the member entity.
 */
export class MemberRepository extends BaseRepository<
  MemberEntity,
  MemberInsertModel,
  MemberSelectModel
> {
  constructor(
    databaseService: DatabaseService<
      MemberEntity,
      MemberInsertModel,
      MemberSelectModel
    >
  ) {
    super(databaseService, MemberTable, MemberEntitySchema)
  }
}
