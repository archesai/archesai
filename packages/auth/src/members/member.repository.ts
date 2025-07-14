import type { DatabaseService } from '@archesai/core'
import type { MemberInsertModel, MemberSelectModel } from '@archesai/database'
import type { MemberEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { MemberTable } from '@archesai/database'
import { MemberEntitySchema } from '@archesai/schemas'

export const createMemberRepository = (
  databaseService: DatabaseService<MemberInsertModel, MemberSelectModel>
) => {
  return createBaseRepository<
    MemberEntity,
    MemberInsertModel,
    MemberSelectModel
  >(databaseService, MemberTable, MemberEntitySchema)
}

export type MemberRepository = ReturnType<typeof createMemberRepository>
