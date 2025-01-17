import { UserEntity } from '@/src/users/entities/user.entity'
import { Injectable, Logger } from '@nestjs/common'
import { PassportStrategy } from '@nestjs/passport'
import { Request } from 'express'
import { ExtractJwt, Strategy } from 'passport-jwt'

import { UsersService } from '../../users/users.service'
import { ConfigService } from '@/src/config/config.service'

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy, 'jwt') {
  private readonly logger = new Logger(JwtStrategy.name)

  constructor(
    private configService: ConfigService,
    private usersService: UsersService
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromExtractors([
        (request: Request) => {
          const bearerToken = ExtractJwt.fromAuthHeaderAsBearerToken()(request)
          if (bearerToken) {
            this.logger.debug(
              `Access Token Extracted From Header: ${bearerToken}`
            )
            return bearerToken
          }
          const cookieToken = request.cookies?.['archesai.accessToken']
          this.logger.debug(
            `Access Token Extracted From Cookies: ${cookieToken}`
          )
          if (cookieToken) {
            return cookieToken
          }
        }
      ]),
      secretOrKey: configService.get('jwt.secret')
    })
  }

  async validate(payload: any): Promise<UserEntity | null> {
    this.logger.debug(
      `Validating JWT token for user: ${JSON.stringify(payload.sub)}`
    )
    if (!payload.sub) {
      return null
    }
    const { sub: id } = payload
    return this.usersService.findOne(id)
  }
}
