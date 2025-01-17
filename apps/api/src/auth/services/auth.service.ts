import { Injectable, Logger, UnauthorizedException } from '@nestjs/common'
import { JwtService } from '@nestjs/jwt'
import { AuthProviderType } from '@prisma/client'
import * as bcrypt from 'bcryptjs'
import { Response } from 'express'

import { UserEntity } from '../../users/entities/user.entity'
import { UsersService } from '../../users/users.service'
import { RegisterDto } from '../dto/register.dto'
import { CookiesDto } from '../dto/token.dto'
import { ConfigService } from '@/src/config/config.service'

@Injectable()
export class AuthService {
  private readonly logger = new Logger(AuthService.name)

  constructor(
    protected jwtService: JwtService,
    protected usersService: UsersService,
    protected configService: ConfigService
  ) {}

  async login(user: UserEntity) {
    this.logger.debug('Logging in user: ' + user.id)
    const accessToken = this.generateAccessToken(user.id)
    const refreshToken = this.generateRefreshToken(user.id)

    // Store refresh token in database
    await this.usersService.setRefreshToken(user.id, refreshToken)

    return {
      accessToken,
      refreshToken
    }
  }

  async refreshAccessToken(refreshToken: string): Promise<CookiesDto> {
    this.logger.debug('Refreshing access token using refresh token')
    const payload = this.jwtService.verify(refreshToken, {
      // secret: this.configService.get("JWT_REFRESH_SECRET"),
    })

    const user = await this.usersService.findOne(payload.sub)

    if (!user || user.refreshToken !== refreshToken) {
      throw new UnauthorizedException('Refresh token is invalid')
    }

    // Generate new tokens
    const newAccessToken = this.generateAccessToken(user.id)
    const newRefreshToken = this.generateRefreshToken(user.id)

    // Update refresh token in the database
    await this.usersService.setRefreshToken(user.id, newRefreshToken)

    return {
      accessToken: newAccessToken,
      refreshToken: newRefreshToken // Return new refresh token
    }
  }

  async register(registerDto: RegisterDto): Promise<UserEntity> {
    this.logger.debug('Registering user: ' + registerDto.email)
    const hashedPassword = await bcrypt.hash(registerDto.password, 10)
    const orgname =
      registerDto.email.split('@')[0] +
      '-' +
      Math.random().toString(36).substring(2, 6)
    const user = await this.usersService.create({
      email: registerDto.email,
      emailVerified: !this.configService.get('email.enabled'),
      password: hashedPassword,
      photoUrl: '',
      username: orgname
    })
    return this.usersService.linkAuthProvider(
      user.id,
      AuthProviderType.LOCAL,
      user.email
    )
  }

  async removeCookies(res: Response) {
    res.clearCookie('archesai.accessToken')
    res.clearCookie('archesai.refreshToken')
  }

  async setCookies(res: Response, authTokens: CookiesDto): Promise<void> {
    res.cookie('archesai.accessToken', authTokens.accessToken, {
      httpOnly: true,
      maxAge: 15 * 60 * 1000, // 15 minutes for access token
      sameSite: 'none',
      secure: true
    })
    res.cookie('archesai.refreshToken', authTokens.refreshToken, {
      httpOnly: true,
      maxAge: 7 * 24 * 60 * 60 * 1000, // 7 days for refresh token
      sameSite: 'none',
      secure: true,
      signed: true
    })
  }

  async getUserFromAccessToken(accessToken: string): Promise<UserEntity> {
    const payload = this.jwtService.verify(accessToken)
    return this.usersService.findOne(payload.sub)
  }

  async verifyToken(token: string): Promise<UserEntity> {
    this.logger.debug('Verifying jwt token: ' + token)
    const { sub: id } = this.jwtService.verify(token)
    return this.usersService.findOne(id)
  }

  private generateAccessToken(id: string): string {
    this.logger.debug('Generating access token for user: ' + id)
    return this.jwtService.sign(
      { sub: id },
      {
        expiresIn: '15m'
      }
    )
  }

  private generateRefreshToken(id: string): string {
    this.logger.debug('Generating refresh token for user: ' + id)
    return this.jwtService.sign(
      { sub: id },
      {
        expiresIn: '7d',
        secret: this.configService.get('jwt.secret') // FIXME Use a different secret for refresh tokens
      }
    )
  }
}
