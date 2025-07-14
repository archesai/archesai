import type { AuthService } from '@archesai/auth'
import type {
  ConfigService,
  EmailService,
  LoggerService,
  RedisService,
  WebsocketsService
} from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { createAuthService } from '@archesai/auth'
// import { StripeService } from '@archesai/billing'
import {
  createConfigService,
  createEmailService,
  createLogger,
  createRedisService,
  createWebsocketsService
} from '@archesai/core'
import { createDrizzleDatabaseService } from '@archesai/database'

export interface Container {
  authService: AuthService
  configService: ConfigService
  databaseService: DrizzleDatabaseService
  emailService: EmailService
  loggerService: LoggerService
  redisService: RedisService
  // stripeService: StripeService
  websocketsService: WebsocketsService
}

export function createContainer(): Container {
  const configService = createConfigService()
  const loggerService = createLogger(configService)
  const databaseService = createDrizzleDatabaseService(
    configService.get('database.url')
  )
  const emailService = createEmailService(configService)
  const redisService = createRedisService(configService, loggerService.logger)
  const authService = createAuthService(databaseService)
  const websocketsService = createWebsocketsService(
    configService,
    redisService,
    loggerService.logger
  )
  // const stripeService = new StripeService(configService)

  return {
    authService,
    configService,
    databaseService,
    emailService,
    loggerService,
    redisService,
    // stripeService,
    websocketsService
  }
}
