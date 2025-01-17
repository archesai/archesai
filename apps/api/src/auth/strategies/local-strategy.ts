import { UserEntity } from '@/src/users/entities/user.entity'
import { Injectable, Logger } from '@nestjs/common'
import { UnauthorizedException } from '@nestjs/common'
import { PassportStrategy } from '@nestjs/passport'
import * as bcrypt from 'bcryptjs'
import { Strategy } from 'passport-local'

import { UsersService } from '../../users/users.service'

@Injectable()
export class LocalStrategy extends PassportStrategy(Strategy) {
  private readonly logger = new Logger(LocalStrategy.name)

  constructor(private usersService: UsersService) {
    super({ usernameField: 'email' })
  }

  async validate(email: string, password: string): Promise<UserEntity> {
    try {
      const user = await this.usersService.findOneByEmail(email)
      const match = await bcrypt.compare(password, user.password)
      if (match) {
        return user
      } else {
        throw new Error()
      }
    } catch (error: any) {
      this.logger.debug(
        `Invalid credentials for user with e-mail ${email}:` + error.message
      )
      throw new UnauthorizedException('Invalid credentials')
    }
  }
}
