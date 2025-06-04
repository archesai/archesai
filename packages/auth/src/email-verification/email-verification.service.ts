import { getEmailVerificationHtml, Logger } from '@archesai/core'

import type { CreateEmailVerificationRequest } from '#email-verification/dto/create-email-verification-request.dto'
import type { UpdateEmailVerificationRequest } from '#email-verification/dto/update-email-verification-request.dto'
import type { UsersService } from '#users/users.service'
import type { VerificationTokensService } from '#verification-tokens/verification-tokens.service'

/**
 * Service for managing email verifications.
 */
export class EmailVerificationService {
  private readonly logger = new Logger(EmailVerificationService.name)
  private readonly usersService: UsersService
  private readonly verificationTokensService: VerificationTokensService

  constructor(
    usersService: UsersService,
    verificationTokensService: VerificationTokensService
  ) {
    this.usersService = usersService
    this.verificationTokensService = verificationTokensService
  }

  public async confirm(
    updateEmailVerificationRequest: UpdateEmailVerificationRequest
  ): Promise<void> {
    this.logger.debug('req', { updateEmailVerificationRequest })
    const { userId } = await this.verificationTokensService.verify(
      'EMAIL_VERIFICATION',
      updateEmailVerificationRequest.token
    )
    await this.usersService.update(userId, {
      emailVerified: new Date().toISOString()
    })
  }

  public async request(
    createEmailVerificationRequest: CreateEmailVerificationRequest
  ): Promise<void> {
    this.logger.debug('req', { createEmailVerificationRequest })
    const token = await this.verificationTokensService.create(
      'EMAIL_VERIFICATION',
      createEmailVerificationRequest.userId,
      24
    )
    return this.verificationTokensService.sendNotification(
      'EMAIL_VERIFICATION',
      createEmailVerificationRequest.email,
      token,
      'Verify Your Email Address',
      getEmailVerificationHtml
    )
  }
}
