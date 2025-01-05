import { PartialType, PickType } from '@nestjs/swagger'

import { UserEntity } from '../entities/user.entity'

export class UpdateUserDto extends PartialType(
  PickType(UserEntity, ['defaultOrgname', 'firstName', 'lastName'] as const)
) {}
