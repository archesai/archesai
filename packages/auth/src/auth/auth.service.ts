import type { ArchesApiRequest, ArchesApiResponse } from '@archesai/core'
import type { AccessTokenEntity } from '@archesai/domain'

import { Logger } from '@archesai/core'

import type { AccessTokensService } from '#access-tokens/access-tokens.service'

/**
 * Service for managing authentication.
 */
export class AuthenticationService {
  private readonly accessTokensService: AccessTokensService
  private readonly logger = new Logger(AuthenticationService.name)

  constructor(accessTokensService: AccessTokensService) {
    this.accessTokensService = accessTokensService
  }

  public async login(userId: string, res?: ArchesApiResponse): Promise<void> {
    this.logger.debug('attempting to login', { userId })
    const accessTokens = await this.accessTokensService.create(userId)
    if (res) {
      this.logger.debug('request was passed, setting cookies')
      this.setCookies(res, accessTokens)
    } else {
      this.logger.debug('request was not passed, not setting cookies')
    }
  }

  public async logout(
    req?: ArchesApiRequest,
    res?: ArchesApiResponse
  ): Promise<void> {
    if (res) {
      this.logger.debug('response was passed, removing cookies')
      this.removeCookies(res)
      this.logger.debug('deleted cookies')
    } else {
      this.logger.debug('response was not passed, not removing cookies')
    }
    if (req) {
      this.logger.debug('request was passed, deleting cookies')
      await req.logOut()
    } else {
      this.logger.debug('request was not passed, not deleting cookies')
    }
  }

  private removeCookies(res: ArchesApiResponse) {
    res.clearCookie('archesai.accessToken')
    res.clearCookie('archesai.refreshToken')
    this.logger.debug('removed cookies in response')
  }

  private setCookies(
    res: ArchesApiResponse,
    accessTokens: AccessTokenEntity
  ): void {
    res.cookie('archesai.accessToken', accessTokens.accessToken, {
      httpOnly: true,
      maxAge: 15 * 60 * 1000, // 15 minutes for access token
      sameSite: 'none',
      secure: true
    })
    res.cookie('archesai.refreshToken', accessTokens.refreshToken, {
      httpOnly: true,
      maxAge: 7 * 24 * 60 * 60 * 1000, // 7 days for refresh token
      sameSite: 'none',
      secure: true,
      signed: true
    })
    this.logger.debug('set cookies in response')
  }
}
