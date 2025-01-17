import { HttpModule } from '@nestjs/axios'
import { BullModule } from '@nestjs/bullmq'
import { Module } from '@nestjs/common'
import { APP_GUARD } from '@nestjs/core'
import { JwtModule } from '@nestjs/jwt'
import { MulterModule } from '@nestjs/platform-express'
import { ScheduleModule } from '@nestjs/schedule'
import { readFileSync } from 'fs'

import { ApiTokensModule } from '@/src/api-tokens/api-tokens.module'
import { AudioModule } from '@/src/audio/audio.module'
import { AuthModule } from '@/src/auth/auth.module'
import { DeactivatedGuard } from '@/src/auth/guards/deactivated.guard'
import { MembershipGuard } from '@/src/auth/guards/membership.guard'
import { BillingModule } from '@/src/billing/billing.module'
import { CommonModule } from '@/src/common/common.module'
import { ContentModule } from '@/src/content/content.module'
import { EmailModule } from '@/src/email/email.module'
import { LabelsModule } from '@/src/labels/labels.module'
import { LlmModule } from '@/src/llm/llm.module'
import { MembersModule } from '@/src/members/members.module'
import { OrganizationsModule } from '@/src/organizations/organizations.module'
import { PipelinesModule } from '@/src/pipelines/pipelines.module'
import { PrismaModule } from '@/src/prisma/prisma.module'
import { RunpodModule } from '@/src/runpod/runpod.module'
import { RunsModule } from '@/src/runs/runs.module'
import { SpeechModule } from '@/src/speech/speech.module'
import { StorageModule } from '@/src/storage/storage.module'
import { ToolsModule } from '@/src/tools/tools.module'
import { UsersModule } from '@/src/users/users.module'
import { WebsocketsModule } from '@/src/websockets/websockets.module'
import { ConfigModule } from '@/src/config/config.module'
import { ApiTokenRestrictedDomainGuard } from '@/src/auth/guards/api-token-restricted-domain.guard'
import { HealthModule } from '@/src/health/health.module'
import { ScraperModule } from '@/src/scraper/scraper.module'
import { ConfigService } from '@/src/config/config.service'
import { APP_INTERCEPTOR } from '@nestjs/core'
import { LoggerErrorInterceptor, LoggerModule, Params } from 'nestjs-pino'

@Module({
  imports: [
    LoggerModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => {
        const targets = []
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
        targets.push({
          options: {
            singleLine: true,
            colorize: true
          },
          target: 'pino-pretty'
        })

        const loggerConfig: Params = {
          pinoHttp: {
            customProps: (req: any) => ({
              context: 'HTTP',
              origin: req?.headers?.origin
            }),
            messageKey: 'message',
            customErrorMessage() {
              return 'http request error'
            },
            customReceivedMessage() {
              return 'http request received'
            },
            customSuccessMessage() {
              return 'http request successful'
            },
            formatters: configService.get('logging.gcpfix')
              ? {
                  level(label: string) {
                    return { level: label, severity: label.toUpperCase() }
                  }
                }
              : undefined,
            level: configService.get('logging.level'),
            redact: {
              paths: ['req', 'res'],
              remove: true
            },
            transport: targets.length > 0 ? { targets } : undefined
          }
        }
        return loggerConfig
      }
    }),
    ConfigModule,
    CommonModule,
    PipelinesModule,
    AuthModule,
    UsersModule,
    OrganizationsModule,
    MembersModule,
    BullModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        connection: {
          host: configService.get('redis.host'),
          password: configService.get('redis.auth'),
          port: configService.get('redis.port'),
          tls: configService.get('redis.ca')
            ? {
                ca: readFileSync(configService.get('redis.ca')!),
                rejectUnauthorized: false
              }
            : undefined
        }
      })
    }),
    JwtModule,
    BillingModule,
    HttpModule,
    PrismaModule,
    EmailModule,
    MulterModule,
    ApiTokensModule,
    LlmModule,
    LabelsModule,
    StorageModule.forRoot(),
    WebsocketsModule,
    AudioModule,
    ScheduleModule.forRoot(),
    ContentModule,
    RunpodModule,
    SpeechModule,
    ToolsModule,
    RunsModule,
    HealthModule,
    ScraperModule
  ],
  providers: [
    {
      provide: APP_GUARD,
      useClass: DeactivatedGuard
    },
    {
      provide: APP_GUARD,
      useClass: MembershipGuard
    },
    {
      provide: APP_GUARD,
      useClass: ApiTokenRestrictedDomainGuard
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: LoggerErrorInterceptor
    }
  ]
})
export class AppModule {}
