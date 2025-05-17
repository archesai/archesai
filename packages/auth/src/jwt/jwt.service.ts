/* eslint-disable @typescript-eslint/no-unnecessary-type-parameters */
import type { DecodeOptions, SignOptions, VerifyOptions } from 'jsonwebtoken'

import jwt from 'jsonwebtoken'

export interface JwtServiceOptions {
  secret: string
  signOptions?: SignOptions | undefined
  verifyOptions?: undefined | VerifyOptions
}

/**
 * Service for creating, verifying, and decoding JWTs.
 */
export class JwtService {
  private readonly secret: string
  private readonly signOptions?: SignOptions | undefined
  private readonly verifyOptions?: undefined | VerifyOptions

  constructor(options: JwtServiceOptions) {
    this.secret = options.secret
    this.signOptions = options.signOptions
    this.verifyOptions = options.verifyOptions
  }

  /**
   * Decode a JWT without verifying its signature.
   * @param token The JWT string to decode.
   * @param options Optional decode options.
   * @returns The decoded payload or `null` if invalid.
   */
  public decode<T extends object>(
    token: string,
    options?: DecodeOptions
  ): null | T {
    return jwt.decode(token, options) as null | T
  }

  /**
   * Sign a payload into a JWT.
   * @param payload The payload to encode.
   * @param options Optional signing options to override defaults.
   * @returns A signed JWT as a string.
   */
  public sign(
    payload: Buffer | object | string,
    options?: SignOptions
  ): string {
    return jwt.sign(payload, this.secret, { ...this.signOptions, ...options })
  }

  /**
   * Verify and decode a JWT.
   * @param token The JWT string to verify.
   * @param options Optional verify options to override defaults.
   * @returns The decoded payload if valid.
   * @throws If the token is invalid or expired.
   */
  public verify<T extends object>(token: string, options?: VerifyOptions): T {
    return jwt.verify(token, this.secret, {
      ...this.verifyOptions,
      ...options
    }) as T
  }
}
