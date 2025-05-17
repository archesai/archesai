import { pino } from 'pino'

import type { ILogger } from '#logging/interface/logger.interface'

import { PinoLoggerAdapter } from '#logging/adapters/pino-logger-adapter'

/**
 * Service for configuring and getting the logger instance.
 */
export class LoggerService {
  private static instance?: ILogger

  constructor(options?: pino.LoggerOptions) {
    LoggerService.configure(options)
  }

  public static configure(options?: pino.LoggerOptions) {
    LoggerService.instance = new PinoLoggerAdapter(
      pino(
        options ?? {
          level: 'info',
          transport: {
            target: 'pino-pretty'
          }
        }
      )
    )
  }

  public static getInstance(): ILogger {
    if (!LoggerService.instance) {
      throw new Error(
        'LoggerService not configured. Call LoggerService.configure() first.'
      )
    }
    return LoggerService.instance
  }
}
