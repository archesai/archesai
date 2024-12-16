import { PickType } from '@nestjs/swagger'
import { IsString, MinLength } from 'class-validator'

import { UserEntity } from '@/src//users/entities/user.entity'
import { Expose } from 'class-transformer'

export class RegisterDto extends PickType(UserEntity, ['email']) {
  /**
   * The password to create and/or login to an account
   * @example 'password'
   */
  @MinLength(7)
  @IsString()
  @Expose()
  password: string
}
