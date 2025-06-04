import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { MemberTable } from '@archesai/database'
import { MemberEntity } from '@archesai/domain'

/**
 * Repository for interacting with the member entity.
 */
export class MemberRepository extends BaseRepository<MemberEntity> {
  constructor(databaseService: DatabaseService<MemberEntity>) {
    super(databaseService, MemberTable, MemberEntity)
  }
}
