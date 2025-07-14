import type {
  ConfigService,
  DatabaseService,
  LoggerService,
  RedisService,
  WebsocketsService
} from '@archesai/core'

import {
  createConfigService,
  createLogger,
  createRedisService,
  createWebsocketsService
} from '@archesai/core'
import { createDrizzleDatabaseService } from '@archesai/database'

export interface Container {
  configService: ConfigService
  databaseService: DatabaseService
  loggerService: LoggerService
  redisService: RedisService
  websocketsService: WebsocketsService
}

export function createContainer(): Container {
  const configService = createConfigService()
  const loggerService = createLogger(configService)
  const databaseService = createDrizzleDatabaseService(
    configService.get('database.url')
  )
  const redisService = createRedisService(configService, loggerService.logger)
  const websocketsService = createWebsocketsService(
    configService,
    redisService,
    loggerService.logger
  )

  return {
    configService,
    databaseService,
    loggerService,
    redisService,
    websocketsService
  }
}
