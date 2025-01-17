import { Body, Controller, Post, Req, Res } from '@nestjs/common'
import { UseGuards } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'
import { Request, Response } from 'express'

import { UserEntity } from '@/src/users/entities/user.entity'
import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { LoginDto } from '@/src/auth/dto/login.dto'
import { RegisterDto } from '@/src/auth/dto/register.dto'
import { LocalAuthGuard } from '@/src/auth/guards/local-auth.guard'
import { AuthService } from '@/src/auth/services/auth.service'

@ApiTags('Authentication')
@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  /**
   * Register a new user
   * @remarks This endpoint will register a new account and return a JWT token which should be provided in your auth headers
   * @throws {409} ConflictException
   */
  @Post('register')
  async register(
    @Body() registerDto: RegisterDto,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<UserEntity> {
    const user = await this.authService.register(registerDto)
    const cookies = await this.authService.login(user)
    await this.authService.setCookies(res, cookies)
    return this.authService.getUserFromAccessToken(cookies.accessToken)
  }

  /**
   * Login with e-mail and password
   * @remarks This endpoint will log you in with your e-mail and password
   * @throws {401} UnauthorizedException
   * @throws {400} BadRequestException
   */
  @UseGuards(LocalAuthGuard)
  @Post('login')
  async login(
    @Body() loginDto: LoginDto,
    @CurrentUser() currentUserDto: UserEntity,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<UserEntity> {
    const cookies = await this.authService.login(currentUserDto)
    await this.authService.setCookies(res, cookies)
    return this.authService.getUserFromAccessToken(cookies.accessToken)
  }

  /**
   * Refresh access token
   * @remarks This endpoint will refresh your access token
   * @throws {401} UnauthorizedException
   */
  @Post('refresh-token')
  async refreshToken(
    @Req() req: Request,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<UserEntity> {
    const refreshToken = req?.signedCookies?.['archesai.refreshToken']
    const cookies = await this.authService.refreshAccessToken(refreshToken)
    await this.authService.setCookies(res, cookies)
    return this.authService.getUserFromAccessToken(cookies.accessToken)
  }

  /**
   * Logout of the current session
   * @remarks This endpoint will log you out of the current session
   * @throws {401} UnauthorizedException
   */
  @Post('logout')
  async logout(
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<void> {
    await this.authService.removeCookies(res)
  }
}
