import type { Logger as PinoLogger } from 'pino'

import { pino } from 'pino'

import type { ILogger } from '#logging/interface/logger.interface'

import { PinoLoggerAdapter } from '#logging/adapters/pino-logger-adapter'

/**
 * Service for configuring and getting the logger instance.
 */
export class LoggerService {
  public static readonly defaultOptions = {
    level: 'info',
    messageKey: 'message',
    transport: {
      options: {
        colorize: true,
        messageKey: 'message'
        // singleLine: true
      },
      target: 'pino-pretty'
    }
  }
  private static instance?: ILogger
  private static pinoLogger?: PinoLogger

  constructor(options?: pino.LoggerOptions) {
    LoggerService.configure(options)
  }

  public static configure(options?: pino.LoggerOptions) {
    const pinoLogger = pino(options ?? this.defaultOptions)
    LoggerService.instance = new PinoLoggerAdapter(pinoLogger)
    LoggerService.pinoLogger = pinoLogger
  }

  public static getInstance(): ILogger {
    if (!LoggerService.instance) {
      throw new Error(
        'LoggerService not configured. Call LoggerService.configure() first.'
      )
    }
    return LoggerService.instance
  }

  public static getPinoInstance(): PinoLogger {
    if (!LoggerService.pinoLogger) {
      throw new Error(
        'LoggerService not configured. Call LoggerService.configure() first.'
      )
    }
    return LoggerService.pinoLogger
  }
}
