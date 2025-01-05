import { HttpModule } from '@nestjs/axios'
import { BullModule } from '@nestjs/bullmq'
import { Module } from '@nestjs/common'
import { APP_GUARD, APP_INTERCEPTOR } from '@nestjs/core'
import { JwtModule } from '@nestjs/jwt'
import { MulterModule } from '@nestjs/platform-express'
import { ScheduleModule } from '@nestjs/schedule'
import { readFileSync } from 'fs'
import { LoggerErrorInterceptor, LoggerModule, Params } from 'nestjs-pino'

import { ApiTokensModule } from './api-tokens/api-tokens.module'
import { AudioModule } from './audio/audio.module'
import { AuthModule } from './auth/auth.module'
import { AppAuthGuard } from './auth/guards/app-auth.guard'
import { DeactivatedGuard } from './auth/guards/deactivated.guard'
import { MembershipGuard } from './auth/guards/organization-role.guard'
import { BillingModule } from './billing/billing.module'
import { CommonModule } from './common/common.module'
import { ContentModule } from './content/content.module'
import { EmailModule } from './email/email.module'
import { EmbeddingsModule } from './embeddings/embeddings.module'
import { LabelsModule } from './labels/labels.module'
import { LLMModule } from './llm/llm.module'
import { MembersModule } from './members/members.module'
import { OrganizationsModule } from './organizations/organizations.module'
import { PipelinesModule } from './pipelines/pipelines.module'
import { PrismaModule } from './prisma/prisma.module'
import { RunpodModule } from './runpod/runpod.module'
import { RunsModule } from './runs/runs.module'
import { SpeechModule } from './speech/speech.module'
import { StorageModule } from './storage/storage.module'
import { ToolsModule } from './tools/tools.module'
import { UsersModule } from './users/users.module'
import { WebsocketsGateway } from './websockets/websockets.gateway'
import { WebsocketsModule } from './websockets/websockets.module'
import { ArchesConfigModule } from './config/config.module'
import { ApiTokenRestrictedDomainGuard } from './auth/guards/api-token-restricted-domain.guard'
import { HealthModule } from './health/health.module'
import { ScraperModule } from './scraper/scraper.module'
import { ArchesConfigService } from './config/config.service'

@Module({
  controllers: [],
  imports: [
    CommonModule,
    PipelinesModule,
    LoggerModule.forRootAsync({
      imports: [ArchesConfigModule],
      inject: [ArchesConfigService],
      useFactory: (configService: ArchesConfigService) => {
        const loggerConfig: Params = {
          pinoHttp: {
            customProps: (req: any, res) => ({
              context: 'HTTP',
              origin: req?.headers?.origin,
              params: req?.params,
              path: req?.path?.split('?')[0],
              query: req?.query,
              statusCode: res?.statusCode
            }),
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
            transport: {
              targets: [
                ...(configService.get('logging.loki.enabled')
                  ? [
                      {
                        options: {
                          host: configService.get('logging.loki.host'),
                          json: true,
                          labels: {
                            app: 'archesai',
                            environment: 'production'
                          }
                        },
                        target: 'pino-loki'
                      }
                    ]
                  : []),
                {
                  options: {
                    colorize: true,
                    singleLine: true
                  },
                  target: 'pino-pretty'
                }
              ]
            }
          }
        }
        return loggerConfig
      }
    }),
    AuthModule,
    UsersModule,
    OrganizationsModule,
    MembersModule,
    ArchesConfigModule,
    BullModule.forRootAsync({
      imports: [ArchesConfigModule],
      inject: [ArchesConfigService],
      useFactory: async (configService: ArchesConfigService) => ({
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
    EmbeddingsModule,
    LLMModule,
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
    WebsocketsGateway,
    {
      provide: APP_GUARD,
      useClass: AppAuthGuard
    },
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
