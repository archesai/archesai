// src/email-change/email-change.service.ts

import {
  BadRequestException,
  ConflictException,
  Injectable,
  Logger
} from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { ARTokenType } from '@prisma/client'

import { EmailService } from '../../email/email.service'
import { getEmailChangeConfirmationHtml } from '../../email/templates'
import { UsersService } from '../../users/users.service'
import { ConfirmationTokenDto } from '../dto/confirmation-token.dto'
import { EmailRequestDto } from '../dto/email-request.dto'
import { ARTokensService } from './ar-tokens.service' // Import TokenService
import { AuthService } from './auth.service'

@Injectable()
export class EmailChangeService {
  private readonly logger = new Logger(EmailChangeService.name)

  constructor(
    private readonly emailService: EmailService,
    private readonly usersService: UsersService,
    private readonly configService: ConfigService,
    private readonly arTokensService: ARTokensService,
    private readonly authService: AuthService
  ) {}

  async confirm(confirmationTokenDto: ConfirmationTokenDto) {
    const { newEmail, userId } = await this.arTokensService.verifyToken(
      ARTokenType.EMAIL_CHANGE,
      confirmationTokenDto.token
    )
    if (!newEmail) {
      throw new BadRequestException('New email is missing.')
    }

    const user = await this.usersService.setEmail(userId, newEmail)
    return this.authService.login(user)
  }

  async request(
    userId: string,
    emailRequestDto: EmailRequestDto
  ): Promise<void> {
    const user = await this.usersService.findOne(userId)
    let newEmailInUse = false
    try {
      await this.usersService.findOneByEmail(emailRequestDto.email)
      newEmailInUse = true
    } catch (error) {
      this.logger.warn(error)
    }
    if (newEmailInUse) {
      throw new ConflictException('New email is already in use.')
    }

    // Generate an email change token (expires in 24 hours) with additional data
    const token = await this.arTokensService.createToken(
      ARTokenType.EMAIL_CHANGE,
      user.id,
      24, // 24 hours expiry
      { newEmail: emailRequestDto.email }
    )

    // Create an email change confirmation link containing the token
    const changeEmailLink = `${this.configService.get('FRONTEND_HOST')}/confirm?type=email-change&token=${token}`

    // Generate the HTML content for the email
    const htmlContent = getEmailChangeConfirmationHtml(
      changeEmailLink,
      user.email
    )

    // Send confirmation email to the new email address
    await this.emailService.sendMail({
      html: htmlContent,
      subject: 'Confirm Your Email Change',
      to: emailRequestDto.email
    })
  }
}
