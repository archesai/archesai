import { PickType } from '@nestjs/swagger'

import { MemberEntity } from '../entities/member.entity'

export class CreateMemberDto extends PickType(MemberEntity, ['inviteEmail', 'role'] as const) {}
