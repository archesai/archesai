import { Strategy as PassportTwitterStrategy } from 'passport-twitter'

import type { ConfigService } from '@archesai/core'

import { catchErrorAsync, Logger, NotFoundException } from '@archesai/core'

import type { AccountsService } from '#accounts/accounts.service'
import type { UsersService } from '#users/users.service'

/**
 * Strategy for authenticating with Twitter.
 */
export class TwitterStrategy extends PassportTwitterStrategy {
  private readonly logger = new Logger(TwitterStrategy.name)

  constructor(
    accountsService: AccountsService,
    configService: ConfigService,
    usersService: UsersService
  ) {
    super(
      {
        callbackURL: configService.get('auth.twitter.callbackURL'),
        consumerKey: configService.get('auth.twitter.consumerKey'),
        consumerSecret: configService.get('auth.twitter.consumerSecret'),
        includeEmail: true
      },
      async (_token, _tokenSecret, profile, cb) => {
        const [findAccountError, account] = await catchErrorAsync(
          accountsService.findByProviderAndProviderAccountId(
            'TWITTER',
            profile.id
          )
        )

        if (account) {
          this.logger.debug('found existing account', { account })
          const [findUserError, user] = await catchErrorAsync(
            usersService.findOne(account.userId)
          )
          if (findUserError) {
            this.logger.error('unexpected error finding user', {
              error: findUserError
            })
            cb(findUserError, false)
            return
          }
          cb(null, user)
          return
        }

        if (!(findAccountError instanceof NotFoundException)) {
          this.logger.error('unexpected error finding account', {
            error: findAccountError
          })
          cb(findAccountError, false)
          return
        }

        this.logger.debug('account not found, creating new one')
        const email = (profile.emails ?? [])[0]?.value
        const username = profile.username
        if (!email) {
          cb(new Error('No email found'), false)
          return
        }

        // Create a new user
        const [createUserError, user] = await catchErrorAsync(
          usersService.create({
            deactivated: false,
            email,
            emailVerified: new Date().toISOString(),
            name: profile.displayName,
            orgname: username,
            ...(profile.photos && profile.photos.length > 0 ?
              { image: profile.photos[0]!.value }
            : {})
          })
        )
        if (createUserError) {
          this.logger.error('unexpected error creating user', {
            error: createUserError
          })
          cb(createUserError, false)
          return
        }

        // Create a new account
        const [createAccountError, newAccount] = await catchErrorAsync(
          accountsService.create({
            authType: 'oauth',
            provider: 'LOCAL',
            providerAccountId: profile.id,
            userId: user.id
          })
        )
        if (createAccountError) {
          this.logger.error('unexpected error creating account', {
            error: createAccountError
          })
          cb(createAccountError, false)
          return
        }

        this.logger.debug('created account for user', { newAccount, user })
        cb(null, user)
      }
    )
  }
}
