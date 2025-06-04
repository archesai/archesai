import { ExtractJwt, Strategy as PassportJwtStrategy } from 'passport-jwt'

import type { ArchesApiRequest, ConfigService } from '@archesai/core'
import type {
  AccessTokenDecodedJwt,
  ApiTokenDecodedJwt
} from '@archesai/domain'

import { Logger, UnauthorizedException } from '@archesai/core'

import type { AccountsService } from '#accounts/accounts.service'
import type { UsersService } from '#users/users.service'

/**
 * Strategy for authenticating with JWT.
 */
export class JwtStrategy extends PassportJwtStrategy {
  private readonly logger = new Logger(JwtStrategy.name)

  constructor(
    accountsService: AccountsService,
    configService: ConfigService,
    usersService: UsersService
  ) {
    super(
      {
        jwtFromRequest: ExtractJwt.fromExtractors([
          (request: ArchesApiRequest) => {
            const bearerToken =
              ExtractJwt.fromAuthHeaderAsBearerToken()(request)
            if (bearerToken) {
              this.logger.debug(`access token extracted from auth header`)
              return bearerToken
            }
            const cookieToken: unknown = request.cookies['archesai.accessToken']
            if (typeof cookieToken === 'string' && cookieToken.length > 0) {
              this.logger.debug(`access token extracted from cookies`)
              return cookieToken
            }
            this.logger.debug(`access token not found`)
            return ''
          }
        ]),
        secretOrKey: configService.get('jwt.secret')
      },
      async (payload: AccessTokenDecodedJwt | ApiTokenDecodedJwt, done) => {
        try {
          this.logger.debug(`validating jwt`, { payload })
          if (!payload.sub) {
            throw new UnauthorizedException('Missing subject in payload')
          }
          const account =
            await accountsService.findByProviderAndProviderAccountId(
              'LOCAL',
              payload.sub
            )
          const user = await usersService.findOne(account.userId)
          done(null, user)
          return
        } catch (error) {
          done(error, false)
          return
        }
      }
    )
  }
}
