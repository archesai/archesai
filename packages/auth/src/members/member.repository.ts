import type { DatabaseService, MemberSelectModel } from '@archesai/database'
import type { MemberEntity } from '@archesai/schemas'

import { createBaseRepository, MemberTable } from '@archesai/database'
import { MemberEntitySchema } from '@archesai/schemas'

export const createMemberRepository = (databaseService: DatabaseService) => {
  return createBaseRepository<MemberEntity, MemberSelectModel>(
    databaseService,
    MemberTable,
    MemberEntitySchema
  )
}

export type MemberRepository = ReturnType<typeof createMemberRepository>
