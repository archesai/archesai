import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { MemberEntity } from '@archesai/schemas'

import { createBaseRepository, MemberTable } from '@archesai/database'
import { MemberEntitySchema } from '@archesai/schemas'

export const createMemberRepository = (
  databaseService: DatabaseService
): BaseRepository<
  MemberEntity,
  (typeof MemberTable)['$inferInsert'],
  (typeof MemberTable)['$inferSelect']
> => {
  return createBaseRepository<MemberEntity>(
    databaseService,
    MemberTable,
    MemberEntitySchema
  )
}

export type MemberRepository = ReturnType<typeof createMemberRepository>
