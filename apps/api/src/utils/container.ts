import type { AuthService } from '@archesai/auth'
import type {
  ConfigService,
  EmailService,
  LoggerService,
  RedisService,
  WebsocketsService
} from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { StorageService } from '@archesai/storage'

import { createAuthService } from '@archesai/auth'
import { StripeService } from '@archesai/billing'
import {
  createConfigService,
  createEmailService,
  createLogger,
  createRedisService,
  createWebsocketsService
} from '@archesai/core'
import { createDatabaseService } from '@archesai/database'
import { createStorageService } from '@archesai/storage'

export interface Container {
  authService: AuthService
  configService: ConfigService
  databaseService: DatabaseService
  emailService: EmailService
  loggerService: LoggerService
  redisService: RedisService
  storageService: StorageService
  stripeService?: StripeService | undefined
  websocketsService: WebsocketsService
}

export function createContainer(): Container {
  const configService = createConfigService()
  const loggerService = createLogger(configService)
  const databaseService = createDatabaseService(
    configService.get('database.url')
  )
  const emailService = createEmailService(configService)
  const redisService = createRedisService(configService, loggerService.logger)
  const authService = createAuthService(databaseService, configService)
  const websocketsService = createWebsocketsService(
    configService,
    redisService,
    loggerService.logger
  )
  const storageService = createStorageService(
    configService,
    loggerService.logger
  )

  let stripeService: StripeService | undefined
  if (configService.get('billing.mode') === 'enabled') {
    stripeService = new StripeService(configService)
  } else {
    loggerService.logger.warn(
      'Stripe service is not initialized because billing mode is disabled.'
    )
  }

  return {
    authService,
    configService,
    databaseService,
    emailService,
    loggerService,
    redisService,
    storageService,
    stripeService,
    websocketsService
  }
}
