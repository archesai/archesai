import type { ModuleMetadata } from '@archesai/core'
import type {
  ApiTokenInsertModel,
  ApiTokenSelectModel
} from '@archesai/database'

import {
  ConfigModule,
  ConfigService,
  createModule,
  DatabaseModule,
  DatabaseService,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { ApiTokenRepository } from '#api-tokens/api-token.repository'
import { ApiTokensController } from '#api-tokens/api-tokens.controller'
import { ApiTokensService } from '#api-tokens/api-tokens.service'
import { JwtModule } from '#jwt/jwt.module'
import { JwtService } from '#jwt/jwt.service'

export const ApiTokensModuleDefinition: ModuleMetadata = {
  exports: [ApiTokensService],
  imports: [ConfigModule, DatabaseModule, JwtModule, WebsocketsModule],
  providers: [
    {
      inject: [ApiTokensService],
      provide: ApiTokensController,
      useFactory: (apiTokensService: ApiTokensService) =>
        new ApiTokensController(apiTokensService)
    },
    {
      inject: [
        ApiTokenRepository,
        ConfigService,
        JwtService,
        WebsocketsService
      ],
      provide: ApiTokensService,
      useFactory: (
        apiTokenRepository: ApiTokenRepository,
        configService: ConfigService,
        jwtService: JwtService,
        websocketsService: WebsocketsService
      ) =>
        new ApiTokensService(
          apiTokenRepository,
          configService,
          jwtService,
          websocketsService
        )
    },
    {
      inject: [DatabaseService],
      provide: ApiTokenRepository,
      useFactory: (
        databaseService: DatabaseService<
          ApiTokenInsertModel,
          ApiTokenSelectModel
        >
      ) => new ApiTokenRepository(databaseService)
    }
  ]
}

export const ApiTokensModule = (() =>
  createModule(class ApiTokensModule {}, ApiTokensModuleDefinition))()
