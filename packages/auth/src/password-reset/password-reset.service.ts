import { getPasswordResetHtml, Logger } from '@archesai/core'

import type { AccountsService } from '#accounts/accounts.service'
import type { HashingService } from '#hashing/hashing.service'
import type { CreatePasswordResetRequest } from '#password-reset/dto/create-password-reset.req.dto'
import type { UpdatePasswordResetRequest } from '#password-reset/dto/update-password-reset.req.dto'
import type { VerificationTokensService } from '#verification-tokens/verification-tokens.service'

/**
 * Service for password reset.
 */
export class PasswordResetService {
  private readonly accountsService: AccountsService
  private readonly hashingService: HashingService
  private readonly logger = new Logger(PasswordResetService.name)
  private readonly verificationTokensService: VerificationTokensService

  constructor(
    accountsService: AccountsService,
    hashingService: HashingService,
    verificationTokensService: VerificationTokensService
  ) {
    this.accountsService = accountsService
    this.hashingService = hashingService
    this.verificationTokensService = verificationTokensService
  }

  public async confirm(
    updatePasswordResetRequest: UpdatePasswordResetRequest
  ): Promise<void> {
    const { userId } = await this.verificationTokensService.verify(
      'PASSWORD_RESET',
      updatePasswordResetRequest.token
    )

    const account =
      await this.accountsService.findByProviderAndProviderAccountId(
        'LOCAL',
        userId
      )

    await this.accountsService.update(account.id, {
      hashed_password: await this.hashingService.hashPassword(
        updatePasswordResetRequest.newPassword
      )
    })
  }

  public async request(
    createPasswordResetRequest: CreatePasswordResetRequest
  ): Promise<void> {
    this.logger.debug('requesting password reset', {
      createPasswordResetRequest
    })
    const { userId } =
      await this.accountsService.findByProviderAndProviderAccountId(
        'LOCAL',
        createPasswordResetRequest.email
      )
    this.logger.debug('found userId', { userId })

    const token = await this.verificationTokensService.create(
      'PASSWORD_RESET',
      userId,
      1
    )
    this.logger.debug('created token, sending email', { token })
    return this.verificationTokensService.sendNotification(
      'PASSWORD_RESET',
      createPasswordResetRequest.email,
      token,
      'Reset Your Password',
      getPasswordResetHtml
    )
  }
}
