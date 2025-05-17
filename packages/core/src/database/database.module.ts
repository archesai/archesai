import type { DynamicModule } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import {
  DATABASE_SERVICE_TOKEN,
  DatabaseHostModule
} from '#database/database-host.module'
import { DatabaseService } from '#database/database.service'
import { Module } from '#utils/nest'

@Module({
  exports: [DatabaseHostModule, DatabaseService],
  imports: [DatabaseHostModule],
  providers: [
    {
      provide: DatabaseService,
      useExisting: DATABASE_SERVICE_TOKEN
    }
  ]
})
export class DatabaseModule {
  public static forRootAsync(
    databaseServiceFactory: (databaseString: string) => DatabaseService
  ): DynamicModule {
    return {
      exports: [DatabaseHostModule, DatabaseService],
      imports: [DatabaseHostModule, ConfigModule],
      module: DatabaseModule,
      providers: [
        {
          inject: [ConfigService],
          provide: DatabaseService,
          useFactory: (configService: ConfigService) =>
            databaseServiceFactory(configService.get('database.url'))
        }
      ]
    }
  }
}
