import type { TransportTargetOptions } from 'pino'

import { pino } from 'pino'

import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { LoggerService } from '#logging/logger.service'
import { Module } from '#utils/nest'

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
            options: {
              colorize: true
              // singleLine: true
            },
            target: 'pino-pretty'
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
          ...(configService.get('logging.gcpfix') ?
            {
              formatters: {
                level: (label: string) => ({
                  level: label,
                  severity: label.toUpperCase()
                })
              }
            }
          : {}),
          transport: {
            targets
          }
        }
        return new LoggerService(loggerConfig)
      }
    }
  ]
}

@Module(LoggingModuleDefinition)
export class LoggingModule {}
