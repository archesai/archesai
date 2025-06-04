import { ExtractJwt, Strategy as PassportJwtStrategy } from 'passport-jwt'

import type { ConfigService } from '@archesai/core'
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
        jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
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
