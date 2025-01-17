import { Body, Controller, Post, Res } from '@nestjs/common'

import { ApiTags } from '@nestjs/swagger'

import { ConfirmationTokenWithNewPasswordDto } from '@/src/auth/dto/confirmation-token-with-new-password.dto'
import { EmailRequestDto } from '@/src/auth/dto/email-request.dto'
import { PasswordResetService } from '@/src/auth/services/password-reset.service'
import { Response } from 'express'
import { AuthService } from '@/src/auth/services/auth.service'
import { UserEntity } from '@/src/users/entities/user.entity'

@ApiTags('Authentication - Password Reset')
@Controller('auth/password-reset')
export class PasswordResetController {
  constructor(
    private passwordResetService: PasswordResetService,
    private authService: AuthService
  ) {}

  /**
   * Confirm password change with a token
   * @remarks This endpoint will confirm your password change with a token
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   */
  @Post('confirm')
  async passwordResetConfirm(
    @Body() confirmPasswordReset: ConfirmationTokenWithNewPasswordDto,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<UserEntity> {
    const cookies =
      await this.passwordResetService.confirm(confirmPasswordReset)
    await this.authService.setCookies(res, cookies)
    return this.authService.getUserFromAccessToken(cookies.accessToken)
  }

  /**
   * Request password reset
   * @remarks This endpoint will request a password reset link
   */
  @Post('request')
  async passwordResetRequest(
    @Body() emailRequestDto: EmailRequestDto
  ): Promise<void> {
    await this.passwordResetService.request(emailRequestDto)
  }
}
