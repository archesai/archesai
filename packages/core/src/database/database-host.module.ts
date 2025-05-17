import type { DatabaseService } from '#database/database.service'
import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { InMemoryDatabaseService } from '#database/adapters/in-memory-database.service'
import { Global, Module } from '#utils/nest'

export const DATABASE_SERVICE_TOKEN = Symbol('DATABASE_SERVICE')

export const DatabaseHostModuleDefinition: ModuleMetadata = {
  exports: [DATABASE_SERVICE_TOKEN],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: DATABASE_SERVICE_TOKEN,
      useFactory: (configService: ConfigService): DatabaseService => {
        const storageType = configService.get('database.type')
        switch (storageType) {
          case 'in-memory':
            return new InMemoryDatabaseService()
          default:
            return new InMemoryDatabaseService()
        }
      }
    }
  ]
}

@Global()
@Module(DatabaseHostModuleDefinition)
export class DatabaseHostModule {}
