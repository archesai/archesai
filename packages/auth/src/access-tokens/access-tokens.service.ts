import type { AccessTokenDecodedJwt } from '@archesai/schemas'

import { Logger, UnauthorizedException } from '@archesai/core'
import { AccessTokenEntity } from '@archesai/schemas'

import type { AccountsService } from '#accounts/accounts.service'
import type { JwtService } from '#jwt/jwt.service'

/**
 * Service for creating and verifying access tokens.
 */
export class AccessTokensService {
  private readonly accountsService: AccountsService
  private readonly jwtService: JwtService
  private readonly logger = new Logger(AccessTokensService.name)

  constructor(accountsService: AccountsService, jwtService: JwtService) {
    this.accountsService = accountsService
    this.jwtService = jwtService
  }

  /**
   * Creates a new access token and refresh token for the given subject (sub).
   * @param sub - The subject identifier for which the tokens are being created.
   * @returns A promise that resolves to an `AccessTokenEntity` containing the generated
   *          access token and refresh token.
   */
  public async create(sub: string): Promise<AccessTokenEntity> {
    this.logger.debug('creating access tokens', { sub })
    const accessToken = this.generate(sub, 'accessToken')
    const refreshToken = this.generate(sub, 'refreshToken')

    this.logger.debug('searching for existing account', { sub })
    const account =
      await this.accountsService.findByProviderAndProviderAccountId(
        'LOCAL',
        sub
      )
    this.logger.debug('got account query result', { account })

    await this.accountsService.update(account.id, {
      refresh_token: refreshToken
    })
    this.logger.debug('updated refresh token in database', {
      refreshToken,
      sub
    })

    return new AccessTokenEntity({
      accessToken,
      refreshToken
    })
  }

  /**
   * Refreshes the access tokens using the provided refresh token.
   * @param refreshToken - The refresh token used to generate new access and refresh tokens.
   * @returns A promise that resolves to an `AccessTokenEntity` containing the new access token
   *          and refresh token.
   * @throws An error if the refresh token is invalid or does not match the stored token.
   */
  public async refresh(refreshToken: string): Promise<AccessTokenEntity> {
    this.logger.debug('refreshing access tokens')
    const payload = this.jwtService.verify<AccessTokenDecodedJwt>(refreshToken)

    const account = await this.accountsService.findOne(payload.sub)
    if (account.refresh_token !== refreshToken) {
      throw new UnauthorizedException('Refresh token is invalid')
    }

    // Generate new tokens
    const newAccessToken = this.generate(payload.sub, 'accessToken')
    const newRefreshToken = this.generate(payload.sub, 'refreshToken')

    // Update refresh token in the database
    await this.accountsService.update(account.id, {
      refresh_token: newRefreshToken
    })

    return new AccessTokenEntity({
      accessToken: newAccessToken,
      refreshToken: newRefreshToken
    })
  }

  /**
   * Verifies the provided access token and decodes its payload.
   * @param accessToken - The JWT access token to be verified.
   * @returns The decoded payload of the access token as an `AccessTokenDecodedJwt` object.
   * @throws {Error} Throws an error if the access token is invalid or verification fails.
   */
  public verify(accessToken: string): AccessTokenDecodedJwt {
    this.logger.debug('verifying access token')
    try {
      return this.jwtService.verify<AccessTokenDecodedJwt>(accessToken)
    } catch (error) {
      this.logger.error('error verifying access token', { error })
      throw new UnauthorizedException('Access token is invalid')
    }
  }

  /**
   * Generates a JSON Web Token (JWT) for the specified subject and token type.
   * @param sub - The subject for which the token is being generated (e.g., user ID).
   * @param type - The type of token to generate, either 'accessToken' or 'refreshToken'.
   *               - 'accessToken': Short-lived token, typically used for authentication.
   *               - 'refreshToken': Long-lived token, used to obtain new access tokens.
   * @returns A signed JWT string.
   */
  private generate(sub: string, type: 'accessToken' | 'refreshToken'): string {
    this.logger.debug(`generating ${type} token`, { sub })
    return this.jwtService.sign(
      { sub },
      {
        expiresIn: type === 'accessToken' ? '15m' : '7d'
      }
    )
  }
}
