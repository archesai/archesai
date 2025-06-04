import { ExtractJwt, Strategy as JwtStrategy } from 'passport-jwt'

import type { ConfigService } from '@archesai/core'
import type { ApiTokenDecodedJwt } from '@archesai/domain'

import { catchErrorAsync, Logger, UnauthorizedException } from '@archesai/core'

import type { AccountsService } from '#accounts/accounts.service'
import type { UsersService } from '#users/users.service'

/**
 * Strategy for authenticating with API key.
 */
export class ApiKeyStrategy extends JwtStrategy {
  private readonly logger = new Logger(ApiKeyStrategy.name)

  constructor(
    accountsService: AccountsService,
    configService: ConfigService,
    usersService: UsersService
  ) {
    super(
      {
        ignoreExpiration: false,
        jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
        secretOrKey: configService.get('jwt.secret')
      },
      async (payload: ApiTokenDecodedJwt, done) => {
        this.logger.debug(`validating api-token`, { payload })
        const [findAccountError, account] = await catchErrorAsync(
          accountsService.findByProviderAndProviderAccountId(
            'API_KEY',
            payload.sub
          )
        )
        if (findAccountError) {
          this.logger.debug(`failed to find account`, {
            error: findAccountError
          })
          done(findAccountError, false)
          return
        }

        this.logger.debug(
          `found user valid for api key, checking memberships`,
          {
            account
          }
        )

        if (!account.scope) {
          this.logger.debug(`account does not have require scope`, {
            account
          })
          done(
            new UnauthorizedException('Account does not have require scope'),
            false
          )
        }

        const [findUserError, user] = await catchErrorAsync(
          usersService.findOne(account.userId)
        )
        if (findUserError) {
          this.logger.debug(`failed to find user`, { error: findUserError })
          done(findUserError, false)
          return
        }

        this.logger.debug(`found user`, { user })
        done(null, user)
      }
    )
  }
}
