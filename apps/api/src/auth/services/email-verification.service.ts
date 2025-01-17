import { Injectable } from '@nestjs/common'
import { ARTokenType } from '@prisma/client'

import { EmailService } from '@/src/email/email.service'
import { getEmailVerificationHtml } from '@/src/email/templates'
import { UsersService } from '@/src/users/users.service'
import { ConfirmationTokenDto } from '@/src/auth/dto/confirmation-token.dto'
import { ARTokensService } from '@/src/auth/services/ar-tokens.service'
import { AuthService } from '@/src/auth/services/auth.service'
import { ConfigService } from '@/src/config/config.service'
import { CookiesDto } from '@/src/auth/dto/token.dto'

@Injectable()
export class EmailVerificationService {
  constructor(
    private readonly usersService: UsersService,
    private readonly emailService: EmailService,
    private readonly configService: ConfigService,
    private readonly authService: AuthService,
    private readonly arTokensService: ARTokensService
  ) {}

  async confirm(
    confirmationTokenDto: ConfirmationTokenDto
  ): Promise<CookiesDto> {
    const { userId } = await this.arTokensService.verifyToken(
      ARTokenType.EMAIL_VERIFICATION,
      confirmationTokenDto.token
    )

    const user = await this.usersService.setEmailVerified(userId)
    return this.authService.login(user)
  }

  async request(userId: string): Promise<void> {
    const user = await this.usersService.findOne(userId)

    const token = await this.arTokensService.createToken(
      ARTokenType.EMAIL_VERIFICATION,
      user.id,
      24
    )

    const verificationLink = `${this.configService.get('frontend.host')}/confirm?type=email-verification&token=${token}`

    const htmlContent = getEmailVerificationHtml(verificationLink)
    await this.emailService.sendMail({
      html: htmlContent,
      subject: 'Verify Your Email Address',
      to: user.email
    })
  }
}
