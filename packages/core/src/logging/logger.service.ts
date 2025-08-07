import type { Logger } from 'pino'
import type { PrettyStream } from 'pino-pretty'

import { pino } from 'pino'
import loki, { pinoLoki as _pinoLoki } from 'pino-loki'
import pretty, { prettyFactory as _pinoPretty } from 'pino-pretty'

import type { ConfigService } from '#config/config.service'

import { PinoLoggerAdapter } from '#logging/adapters/pino-logger-adapter'

export const createLogger = (
  configService: ConfigService
): {
  logger: PinoLoggerAdapter
  pinoLogger: Logger
} => {
  const streams: PrettyStream[] = []

  const prettyTransport = pretty({
    colorize: true,
    messageKey: 'message'
  })
  streams.push(prettyTransport)

  const defaultOptions = {
    level: 'info',
    messageKey: 'message'
  }

  if (configService.get('monitoring.loki.mode') !== 'disabled') {
    const _lokiTransport = loki({
      host: configService.get('monitoring.loki.host'),
      labels: {
        app: 'archesai',
        environment: 'production'
      }
    })
    streams.push(_lokiTransport)
  }
  const loggerConfig: pino.LoggerOptions = {
    level: configService.get('logging.level'),
    messageKey: 'message'
  }
  const pinoLogger = pino({ ...loggerConfig, ...defaultOptions }, ...streams)
  return {
    logger: new PinoLoggerAdapter(pinoLogger),
    pinoLogger
  }
}

export type LoggerService = ReturnType<typeof createLogger>
