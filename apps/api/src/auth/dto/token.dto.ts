import { Expose } from 'class-transformer'
import { IsString } from 'class-validator'

export class TokenDto {
  /**
   * The authorization token that can be used to access Arches AI
   * @example 'supersecretauthorizationtoken'
   */
  @IsString()
  @Expose()
  accessToken: string

  /**
   * The refresh token that can be used to get a new access token
   * @example 'supersecretauthorizationtoken'
   */
  @IsString()
  @Expose()
  refreshToken: string

  constructor({
    accessToken,
    refreshToken
  }: {
    accessToken: string
    refreshToken: string
  }) {
    this.accessToken = accessToken
    this.refreshToken = refreshToken
  }
}
