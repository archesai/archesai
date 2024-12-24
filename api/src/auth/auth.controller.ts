import { Body, Controller, Get, Post, Req, Res } from '@nestjs/common'
import { UseGuards } from '@nestjs/common'
import { AuthGuard } from '@nestjs/passport'
import { ApiBearerAuth, ApiExcludeEndpoint, ApiTags } from '@nestjs/swagger'
import { Request, Response } from 'express'

import { UserEntity } from '../users/entities/user.entity'
import { CurrentUser } from './decorators/current-user.decorator'
import { IsPublic } from './decorators/is-public.decorator'
import { Roles } from './decorators/roles.decorator'
import { ConfirmationTokenWithNewPasswordDto } from './dto/confirmation-token-with-new-password.dto'
import { ConfirmationTokenDto } from './dto/confirmation-token.dto'
import { EmailRequestDto } from './dto/email-request.dto'
import { LoginDto } from './dto/login.dto'
import { RegisterDto } from './dto/register.dto'
import { TokenDto } from './dto/token.dto'
import { LocalAuthGuard } from './guards/local-auth.guard'
import { AuthService } from './services/auth.service'
import { EmailChangeService } from './services/email-change.service'
import { EmailVerificationService } from './services/email-verification.service'
import { PasswordResetService } from './services/password-reset.service'

@ApiTags('Authentication')
@Controller('auth')
export class AuthController {
  constructor(
    private readonly authService: AuthService,
    private passwordResetService: PasswordResetService,
    private emailVerificationService: EmailVerificationService,
    private emailChangeService: EmailChangeService
  ) {}

  /**
   * Confirm e-mail change with a token
   * @remarks This endpoint will confirm your e-mail change with a token
   * @throws {400} BadRequestException
   */
  @IsPublic()
  @Post('/email-change/confirm')
  async emailChangeConfirm(
    @Body() confirmEmailChangeDto: ConfirmationTokenDto
  ) {
    return new TokenDto(
      await this.emailChangeService.confirm(confirmEmailChangeDto)
    )
  }

  /**
   * Request e-mail change with a token
   * @remarks This endpoint will request your e-mail change with a token
   */
  @ApiBearerAuth()
  @Post('/email-change/request')
  async emailChangeRequest(
    @CurrentUser() currentUserDto: UserEntity,
    @Body() emailRequestDto: EmailRequestDto
  ) {
    return this.emailChangeService.request(currentUserDto.id, emailRequestDto)
  }

  /**
   * Confirm e-mail verification with a token
   * @remarks This endpoint will confirm your e-mail with a token
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   */
  @IsPublic()
  @Post('/email-verification/confirm')
  async emailVerificationConfirm(
    @Body() confirmEmailVerificationDto: ConfirmationTokenDto
  ) {
    return new TokenDto(
      await this.emailVerificationService.confirm(confirmEmailVerificationDto)
    )
  }

  /**
   * Request e-mail verification
   * @remarks This endpoint will send an e-mail verification link to you. ADMIN ONLY.
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {403} ForbiddenException
   */
  @ApiBearerAuth()
  @Post('/email-verification/request')
  @Roles('ADMIN')
  async emailVerificationRequest(@CurrentUser() user: UserEntity) {
    return this.emailVerificationService.request(user.id)
  }

  /**
   * Login with e-mail and password
   * @remarks This endpoint will log you in with your e-mail and password
   * @throws {401} UnauthorizedException
   * @throws {400} BadRequestException
   */
  @IsPublic()
  @Post('/login')
  @UseGuards(LocalAuthGuard)
  async login(
    @Body() loginDto: LoginDto,
    @CurrentUser() currentUserDto: UserEntity,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<TokenDto> {
    const tokenDto = await this.authService.login(currentUserDto)
    await this.authService.setCookies(res, tokenDto)
    return tokenDto
  }

  /**
   * Logout of the current session
   * @remarks This endpoint will log you out of the current session
   * @throws {401} UnauthorizedException
   */
  @IsPublic()
  @Post('/logout')
  async logout(
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<void> {
    await this.authService.removeCookies(res)
  }

  /**
   * Confirm password change with a token
   * @remarks This endpoint will confirm your password change with a token
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   */
  @IsPublic()
  @Post('/password-reset/confirm')
  async passwordResetConfirm(
    @Body() confirmPasswordReset: ConfirmationTokenWithNewPasswordDto
  ): Promise<TokenDto> {
    return new TokenDto(
      await this.passwordResetService.confirm(confirmPasswordReset)
    )
  }

  /**
   * Request password reset
   * @remarks This endpoint will request a password reset link
   */
  @IsPublic()
  @Post('/password-reset/request')
  async passwordResetRequest(@Body() emailRequestDto: EmailRequestDto) {
    await this.passwordResetService.request(emailRequestDto)
  }

  /**
   * Refresh access token
   * @remarks This endpoint will refresh your access token
   * @throws {401} UnauthorizedException
   */
  @IsPublic()
  @Post('/refresh-token')
  async refreshToken(
    @Body('refreshToken') refreshToken: string,
    @Req() req: Request,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<TokenDto> {
    const cookies = req?.signedCookies?.['archesai.refreshToken']
    const tokens = await this.authService.refreshAccessToken(
      refreshToken || cookies
    )
    await this.authService.setCookies(res, tokens)
    return tokens
  }

  /**
   * Register a new user
   * @remarks This endpoint will register a new account and return a JWT token which should be provided in your auth headers
   * @throws {409} ConflictException
   */
  @IsPublic()
  @Post('/register')
  async register(@Body() registerDto: RegisterDto): Promise<TokenDto> {
    const user = await this.authService.register(registerDto)
    return this.authService.login(user)
  }

  @ApiExcludeEndpoint()
  @Post('firebase/callback')
  @UseGuards(AuthGuard('firebase-auth'))
  async zfirebaseAuthCallback(
    @CurrentUser() currentUserDto: UserEntity
  ): Promise<TokenDto> {
    return this.authService.login(currentUserDto)
  }

  @ApiExcludeEndpoint()
  @Get('twitter')
  @UseGuards(AuthGuard('twitter'))
  async ztwitterAuth() {}

  @ApiExcludeEndpoint()
  @Get('twitter/callback')
  @UseGuards(AuthGuard('twitter'))
  async ztwitterAuthCallback(
    @CurrentUser() currentUserDto: UserEntity
  ): Promise<TokenDto> {
    return this.authService.login(currentUserDto)
  }
}
