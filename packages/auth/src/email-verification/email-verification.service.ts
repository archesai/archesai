import type {
  CreateEmailVerificationDto,
  UpdateEmailVerificationDto
} from '@archesai/schemas'

import { getEmailVerificationHtml, Logger } from '@archesai/core'

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
    updateEmailVerificationDto: UpdateEmailVerificationDto
  ): Promise<void> {
    this.logger.debug('req', { updateEmailVerificationDto })
    const { userId } = await this.verificationTokensService.verify(
      'EMAIL_VERIFICATION',
      updateEmailVerificationDto.token
    )
    await this.usersService.update(userId, {
      emailVerified: new Date().toISOString()
    })
  }

  public async request(
    createEmailVerificationDto: CreateEmailVerificationDto
  ): Promise<void> {
    this.logger.debug('req', { createEmailVerificationDto })
    const token = await this.verificationTokensService.create(
      'EMAIL_VERIFICATION',
      createEmailVerificationDto.userId,
      24
    )
    return this.verificationTokensService.sendNotification(
      'EMAIL_VERIFICATION',
      createEmailVerificationDto.email,
      token,
      'Verify Your Email Address',
      getEmailVerificationHtml
    )
  }
}
