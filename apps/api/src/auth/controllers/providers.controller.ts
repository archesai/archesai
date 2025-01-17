import { Controller, Get, Post, Res } from '@nestjs/common'
import { UseGuards } from '@nestjs/common'
import { AuthGuard } from '@nestjs/passport'
import { ApiExcludeController } from '@nestjs/swagger'

import { UserEntity } from '@/src/users/entities/user.entity'
import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { AuthService } from '@/src/auth/services/auth.service'
import { Response } from 'express'

@ApiExcludeController()
@Controller('auth/providers')
export class ProvidersController {
  constructor(private readonly authService: AuthService) {}

  @Post('firebase')
  @UseGuards(AuthGuard('firebase-auth'))
  async firebase(
    @CurrentUser() currentUserDto: UserEntity,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<void> {
    const cookies = await this.authService.login(currentUserDto)
    await this.authService.setCookies(res, cookies)
  }

  @Get('twitter')
  @UseGuards(AuthGuard('twitter'))
  async twitter(): Promise<void> {}

  @Get('twitter/callback')
  @UseGuards(AuthGuard('twitter'))
  async twitterCallback(
    @CurrentUser() currentUserDto: UserEntity,
    @Res({
      passthrough: true
    })
    res: Response
  ): Promise<void> {
    const cookies = await this.authService.login(currentUserDto)
    await this.authService.setCookies(res, cookies)
  }
}
