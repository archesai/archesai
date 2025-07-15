import type { Logger, TransportTargetOptions } from 'pino'

import { pino } from 'pino'

import type { ConfigService } from '#config/config.service'

import { PinoLoggerAdapter } from '#logging/adapters/pino-logger-adapter'

export const createLogger = (
  configService: ConfigService
): {
  logger: PinoLoggerAdapter
  pinoLogger: Logger
} => {
  const defaultOptions = {
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
  const targets: TransportTargetOptions[] = [
    {
      ...defaultOptions.transport
    }
  ]
  if (configService.get('monitoring.loki.enabled')) {
    targets.push({
      options: {
        host: configService.get('monitoring.loki.host'),
        json: true,
        labels: {
          app: 'archesai',
          environment: 'production'
        }
      },
      target: 'pino-loki'
    })
  }
  const loggerConfig: pino.LoggerOptions = {
    level: configService.get('logging.level'),
    messageKey: 'message',
    transport: {
      targets
    }
  }
  const pinoLogger = pino(loggerConfig)
  return {
    logger: new PinoLoggerAdapter(pinoLogger),
    pinoLogger
  }
}

export type LoggerService = ReturnType<typeof createLogger>
