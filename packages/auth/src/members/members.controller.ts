import type { Controller } from '@archesai/core'
import type { MemberEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import { MEMBER_ENTITY_KEY, MemberEntitySchema } from '@archesai/schemas'

import type { MembersService } from '#members/members.service'

import { CreateMemberRequestSchema } from '#members/dto/create-member.req.dto'
import { UpdateMemberRequestSchema } from '#members/dto/update-member.req.dto'

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
      CreateMemberRequestSchema,
      UpdateMemberRequestSchema,
      membersService
    )
  }
}
