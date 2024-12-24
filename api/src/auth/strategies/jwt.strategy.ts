import { UserEntity } from '@/src/users/entities/user.entity'
import { Injectable, Logger } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { PassportStrategy } from '@nestjs/passport'
import { Request } from 'express'
import { ExtractJwt, Strategy } from 'passport-jwt'

import { UsersService } from '../../users/users.service'

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  private readonly logger: Logger = new Logger(JwtStrategy.name)

  constructor(
    private configService: ConfigService,
    private usersService: UsersService
  ) {
    super({
      ignoreExpiration: false,
      jwtFromRequest: ExtractJwt.fromExtractors([
        (request: Request) => {
          // Check for token in the Authorization header
          const authToken = ExtractJwt.fromAuthHeaderAsBearerToken()(request)
          this.logger.debug(`authToken: ${authToken}`)
          if (authToken) {
            return authToken
          }

          // Check for token in the signed cookie named 'auth_token'
          const cookieToken = request.cookies?.['archesai.accessToken']
          this.logger.debug(`cookieToken: ${cookieToken}`)
          if (cookieToken) {
            return cookieToken
          }

          return null
        }
      ]),
      secretOrKey: configService.get<string>('JWT_API_TOKEN_SECRET')
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
