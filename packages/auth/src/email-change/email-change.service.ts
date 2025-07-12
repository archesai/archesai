import type {
  CreateEmailChangeDto,
  UpdateEmailChangeDto
} from '@archesai/schemas'

import {
  BadRequestException,
  getEmailChangeConfirmationHtml,
  Logger
} from '@archesai/core'

import type { UsersService } from '#users/users.service'
import type { VerificationTokensService } from '#verification-tokens/verification-tokens.service'

/**
 * Service for managing email changes.
 */
export class EmailChangeService {
  private readonly logger = new Logger(EmailChangeService.name)
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
    updateEmailChangeDto: UpdateEmailChangeDto
  ): Promise<void> {
    this.logger.debug('req', { updateEmailChangeDto })
    const { newEmail, userId } = await this.verificationTokensService.verify(
      'EMAIL_CHANGE',
      updateEmailChangeDto.token
    )

    if (!newEmail) {
      throw new BadRequestException('New email is missing.')
    }
    await this.usersService.update(userId, {
      email: newEmail
    })
  }

  public async request(
    createEmailChangeDto: CreateEmailChangeDto
  ): Promise<void> {
    this.logger.debug('req', { createEmailChangeDto })
    const token = await this.verificationTokensService.create(
      'EMAIL_CHANGE',
      createEmailChangeDto.userId,
      24,
      {
        newEmail: createEmailChangeDto.newEmail
      }
    )
    return this.verificationTokensService.sendNotification(
      'EMAIL_CHANGE',
      createEmailChangeDto.newEmail,
      token,
      'Confirm Your Email Change',
      getEmailChangeConfirmationHtml
    )
  }
}
