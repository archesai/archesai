import { PickType } from '@nestjs/swagger'

import { ApiTokenEntity } from '../entities/api-token.entity'

export class CreateApiTokenDto extends PickType(ApiTokenEntity, [
  'role',
  'domains',
  'name'
] as const) {}
