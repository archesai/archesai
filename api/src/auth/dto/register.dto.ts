import { ApiProperty, PickType } from '@nestjs/swagger'
import { MinLength } from 'class-validator'

import { UserEntity } from '../../users/entities/user.entity'

export class RegisterDto extends PickType(UserEntity, ['email']) {
  @ApiProperty({
    description: 'The password to create and/or login to an account',
    example: 'password',
    minLength: 7
  })
  @MinLength(7)
  password: string
}
