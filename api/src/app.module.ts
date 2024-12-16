import { HttpModule } from '@nestjs/axios'
import { BullModule } from '@nestjs/bullmq'
import { Module } from '@nestjs/common'
import { ConfigModule, ConfigService } from '@nestjs/config'
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
import { validationSchema } from './config/schema'
import { ApiTokenRestrictedDomainGuard } from './auth/guards/api-token-restricted-domain.guard'

@Module({
  controllers: [],
  imports: [
    CommonModule,
    PipelinesModule,
    LoggerModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => {
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
            formatters:
              configService.get<string>('NODE_ENV') === 'production'
                ? {
                    level(label: string) {
                      return { level: label, severity: label.toUpperCase() }
                    }
                  }
                : undefined,
            level: configService.get<string>('LOGGING_LEVEL'),
            redact: {
              paths: ['req', 'res'],
              remove: true
            },
            transport: {
              targets: [
                ...(configService.get<string>('LOKI_HOST')
                  ? [
                      {
                        options: {
                          host: configService.get<string>('LOKI_HOST'),
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
    ConfigModule.forRoot({
      ignoreEnvFile: true,
      isGlobal: true,
      validationSchema: validationSchema
    }),
    BullModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        connection: {
          host: configService.get('REDIS_HOST'),
          password: configService.get('REDIS_AUTH'),
          port: Number(configService.get('REDIS_PORT')),
          tls: configService.get('REDIS_CA_CERT_PATH')
            ? {
                ca: readFileSync(configService.get('REDIS_CA_CERT_PATH')),
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
    RunsModule
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
