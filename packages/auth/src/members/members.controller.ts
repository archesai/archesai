import type { Controller } from '@archesai/core'
import type { MemberEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreateMemberDtoSchema,
  MEMBER_ENTITY_KEY,
  MemberEntitySchema,
  UpdateMemberDtoSchema
} from '@archesai/schemas'

import type { MembersService } from '#members/members.service'

/**
 * Controller for handling members.
 */
export class MembersController
  extends BaseController<MemberEntity>
  implements Controller
{
  constructor(membersService: MembersService) {
    super(
      MEMBER_ENTITY_KEY,
      MemberEntitySchema,
      CreateMemberDtoSchema,
      UpdateMemberDtoSchema,
      membersService
    )
  }
}
