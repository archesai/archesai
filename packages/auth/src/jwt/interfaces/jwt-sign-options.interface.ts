import type jwt from 'jsonwebtoken'

import type { ModuleMetadata, Provider, Type } from '@archesai/core'

export type GetSecretKeyResult = Buffer | jwt.Secret | string

export interface JwtModuleAsyncOptions extends Pick<ModuleMetadata, 'imports'> {
  extraProviders?: Provider[]
  global?: boolean
  inject?: Type[]
  useClass?: Type<JwtOptionsFactory>
  useExisting?: Type<JwtOptionsFactory>
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  useFactory: (...args: any[]) => JwtModuleOptions | Promise<JwtModuleOptions>
}

export interface JwtModuleOptions {
  global?: boolean
  privateKey?: jwt.Secret
  publicKey?: Buffer | string
  secret?: Buffer | string
  secretOrKeyProvider?: (
    requestType: JwtSecretRequestType,
    tokenOrPayload: Buffer | object | string,
    options?: jwt.SignOptions | jwt.VerifyOptions
  ) => jwt.Secret | Promise<jwt.Secret>
  signOptions?: jwt.SignOptions
  verifyOptions?: jwt.VerifyOptions
}

export interface JwtOptionsFactory {
  createJwtOptions(): JwtModuleOptions | Promise<JwtModuleOptions>
}

export type JwtSecretRequestType = 'SIGN' | 'VERIFY'

export interface JwtSignOptions extends jwt.SignOptions {
  privateKey?: jwt.Secret
  secret?: Buffer | string
}

export interface JwtVerifyOptions extends jwt.VerifyOptions {
  publicKey?: Buffer | string
  secret?: Buffer | string
}
