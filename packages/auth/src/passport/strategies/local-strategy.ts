import { Strategy as LocalStrategyBase } from 'passport-local'

import { Logger } from '@archesai/core'

import type { AccountsService } from '#accounts/accounts.service'
import type { HashingService } from '#hashing/hashing.service'
import type { UsersService } from '#users/users.service'

/**
 * Strategy for authenticating with local credentials.
 */
export class LocalStrategy extends LocalStrategyBase {
  private readonly logger = new Logger(LocalStrategy.name)

  constructor(
    accountsService: AccountsService,
    hashingService: HashingService,
    usersService: UsersService
  ) {
    super(
      { usernameField: 'email' },
      async (email: string, password: string, done) => {
        try {
          this.logger.debug(`validating local credentials`, { email })

          const account =
            await accountsService.findByProviderAndProviderAccountId(
              'LOCAL',
              email
            )

          const match = await hashingService.verifyPassword(
            password,
            account.hashed_password ?? ''
          )

          if (!match) {
            this.logger.debug(`invalid credentials`, { email })
            done(undefined, false)
          }

          const user = await usersService.findOne(account.userId)

          done(null, user)
          return
        } catch (error) {
          this.logger.debug(`error during local strategy validation`, {
            email,
            error
          })
          done(error, false)
          return
        }
      }
    )
  }
}
