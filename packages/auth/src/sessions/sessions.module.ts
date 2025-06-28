import type { Strategy } from 'passport'

import type { ModuleMetadata } from '@archesai/core'

import {
  ConfigModule,
  ConfigService,
  Module,
  RedisService
} from '@archesai/core'

import { PassportModule } from '#passport/passport.module'
import { ApiKeyStrategy } from '#passport/strategies/api-key.strategy'
import { JwtStrategy } from '#passport/strategies/jwt.strategy'
import { LocalStrategy } from '#passport/strategies/local-strategy'
import { SessionSerializer } from '#sessions/session.serializer'
import { SessionsController } from '#sessions/sessions.controller'
import { SessionsService } from '#sessions/sessions.service'
import { UsersModule } from '#users/users.module'
import { UsersService } from '#users/users.service'

export const SessionsModuleDefinition: ModuleMetadata = {
  exports: [SessionsService],
  imports: [ConfigModule, PassportModule, UsersModule],
  providers: [
    {
      inject: [
        ApiKeyStrategy,
        ConfigService,
        JwtStrategy,
        LocalStrategy,
        RedisService,
        SessionSerializer
      ],
      provide: SessionsService,
      useFactory: (
        apiKeyStrategy: ApiKeyStrategy,
        configService: ConfigService,
        jwtStrategy: JwtStrategy,
        localStrategy: LocalStrategy,
        redisService: RedisService,
        sessionSerializer: SessionSerializer
      ) => {
        const strategies = {
          'api-key-auth': apiKeyStrategy,
          jwt: jwtStrategy,
          local: localStrategy
        } satisfies Record<string, Strategy>
        return new SessionsService(
          configService,
          redisService,
          strategies,
          sessionSerializer
        )
      }
    },
    {
      inject: [UsersService],
      provide: SessionSerializer,
      useFactory: (usersService: UsersService) =>
        new SessionSerializer(usersService)
    },
    {
      inject: [SessionsService],
      provide: SessionsController,
      useFactory: (sessionsService: SessionsService) =>
        new SessionsController(sessionsService)
    }
  ]
}

@Module(SessionsModuleDefinition)
export class SessionsModule {}
