import type { pino, TransportTargetOptions } from 'pino'

import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { LoggerService } from '#logging/logger.service'
import { createModule } from '#utils/nest'

export const LoggingModuleDefinition: ModuleMetadata = {
  exports: [LoggerService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: LoggerService,
      useFactory: (configService: ConfigService) => {
        const targets: TransportTargetOptions[] = [
          {
            ...LoggerService.defaultOptions.transport
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
        return new LoggerService(loggerConfig)
      }
    }
  ]
}

export const LoggingModule = (() =>
  createModule(class LoggingModule {}, LoggingModuleDefinition))()
