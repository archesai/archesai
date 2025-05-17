import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, Module } from '@archesai/core'
import { UsersModule, UsersService } from '@archesai/organizations'

import { AccountsModule } from '#accounts/accounts.module'
import { AccountsService } from '#accounts/accounts.service'
import { HashingModule } from '#hashing/hashing.module'
import { HashingService } from '#hashing/hashing.service'
import { ApiKeyStrategy } from '#passport/strategies/api-key.strategy'
import { JwtStrategy } from '#passport/strategies/jwt.strategy'
import { LocalStrategy } from '#passport/strategies/local-strategy'

export const PassportModuleDefinition: ModuleMetadata = {
  exports: [LocalStrategy, JwtStrategy, ApiKeyStrategy],
  imports: [AccountsModule, ConfigModule, HashingModule, UsersModule],
  providers: [
    {
      inject: [AccountsService, HashingService, UsersService],
      provide: LocalStrategy,
      useFactory: (
        accountsService: AccountsService,
        hashingService: HashingService,
        usersService: UsersService
      ) => new LocalStrategy(accountsService, hashingService, usersService)
    },
    {
      inject: [AccountsService, ConfigService, UsersService],
      provide: JwtStrategy,
      useFactory: (
        accountsService: AccountsService,
        configService: ConfigService,
        usersService: UsersService
      ) => new JwtStrategy(accountsService, configService, usersService)
    },
    {
      inject: [AccountsService, ConfigService, UsersService],
      provide: ApiKeyStrategy,
      useFactory: (
        accountsService: AccountsService,
        configService: ConfigService,
        usersService: UsersService
      ) => new ApiKeyStrategy(accountsService, configService, usersService)
    }
  ]
}

@Module(PassportModuleDefinition)
export class PassportModule {}
