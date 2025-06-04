import {
  BadRequestException,
  getEmailChangeConfirmationHtml,
  Logger
} from '@archesai/core'

import type { CreateEmailChangeRequest } from '#email-change/dto/create-email-change-request.dto'
import type { UpdateEmailChangeRequest } from '#email-change/dto/update-email-change-request.dto'
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
    updateEmailChangeRequest: UpdateEmailChangeRequest
  ): Promise<void> {
    this.logger.debug('req', { updateEmailChangeRequest })
    const { newEmail, userId } = await this.verificationTokensService.verify(
      'EMAIL_CHANGE',
      updateEmailChangeRequest.token
    )

    if (!newEmail) {
      throw new BadRequestException('New email is missing.')
    }
    await this.usersService.update(userId, {
      email: newEmail
    })
  }

  public async request(
    createEmailChangeRequest: CreateEmailChangeRequest
  ): Promise<void> {
    this.logger.debug('req', { createEmailChangeRequest })
    const token = await this.verificationTokensService.create(
      'EMAIL_CHANGE',
      createEmailChangeRequest.userId,
      24,
      {
        newEmail: createEmailChangeRequest.newEmail
      }
    )
    return this.verificationTokensService.sendNotification(
      'EMAIL_CHANGE',
      createEmailChangeRequest.newEmail,
      token,
      'Confirm Your Email Change',
      getEmailChangeConfirmationHtml
    )
  }
}
