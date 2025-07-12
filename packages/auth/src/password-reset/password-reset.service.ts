import type {
  CreatePasswordResetDto,
  UpdatePasswordResetDto
} from '@archesai/schemas'

import { getPasswordResetHtml, Logger } from '@archesai/core'

import type { AccountsService } from '#accounts/accounts.service'
import type { HashingService } from '#hashing/hashing.service'
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
    updatePasswordResetDto: UpdatePasswordResetDto
  ): Promise<void> {
    const { userId } = await this.verificationTokensService.verify(
      'PASSWORD_RESET',
      updatePasswordResetDto.token
    )

    const account =
      await this.accountsService.findByProviderAndProviderAccountId(
        'LOCAL',
        userId
      )

    await this.accountsService.update(account.id, {
      password: await this.hashingService.hashPassword(
        updatePasswordResetDto.newPassword
      )
    })
  }

  public async request(
    createPasswordResetDto: CreatePasswordResetDto
  ): Promise<void> {
    this.logger.debug('requesting password reset', {
      createPasswordResetDto
    })
    const { userId } =
      await this.accountsService.findByProviderAndProviderAccountId(
        'LOCAL',
        createPasswordResetDto.email
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
      createPasswordResetDto.email,
      token,
      'Reset Your Password',
      getPasswordResetHtml
    )
  }
}
