import { Body, Controller, Post, Res } from '@nestjs/common'

import { ApiTags } from '@nestjs/swagger'

import { UserEntity } from '@/src/users/entities/user.entity'
import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { ConfirmationTokenDto } from '@/src/auth/dto/confirmation-token.dto'

import { EmailVerificationService } from '@/src/auth/services/email-verification.service'
import { AuthService } from '../services/auth.service'
import { Response } from 'express'
import { Authenticated } from '../decorators/authenticated.decorator'

@ApiTags('Authentication - Email Verification')
@Controller('auth/email-verification')
export class EmailVerificationController {
  constructor(
    private emailVerificationService: EmailVerificationService,
    private authService: AuthService
  ) {}

  /**
   * Confirm e-mail verification with a token
   * @remarks This endpoint will confirm your e-mail with a token
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   */
  @Post('confirm')
  async emailVerificationConfirm(
    @Body() confirmEmailVerificationDto: ConfirmationTokenDto,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<UserEntity> {
    const cookies = await this.emailVerificationService.confirm(
      confirmEmailVerificationDto
    )
    await this.authService.setCookies(res, cookies)
    return this.authService.getUserFromAccessToken(cookies.accessToken)
  }

  /**
   * Request e-mail verification
   * @remarks This endpoint will send an e-mail verification link to you. ADMIN ONLY.
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {403} ForbiddenException
   */
  @Authenticated()
  @Post('request')
  async emailVerificationRequest(
    @CurrentUser() user: UserEntity
  ): Promise<void> {
    return this.emailVerificationService.request(user.id)
  }
}
