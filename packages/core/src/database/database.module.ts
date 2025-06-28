import type { DynamicModule } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { DatabaseService } from '#database/database.service'
import { createModule } from '#utils/nest'

export const DATABASE_SERVICE_TOKEN = Symbol('DATABASE_SERVICE')

const DatabaseModuleDefinition = {
  exports: [DatabaseService],
  providers: [
    {
      provide: DatabaseService,
      useExisting: DATABASE_SERVICE_TOKEN
    }
  ]
}

export const DatabaseModule = (() =>
  createModule(
    class DynamicDatabaseModule {
      public static forRootAsync(
        databaseServiceFactory: (databaseString: string) => DatabaseService
      ): DynamicModule {
        return {
          exports: [DATABASE_SERVICE_TOKEN],
          global: true,
          imports: [ConfigModule],
          module: DynamicDatabaseModule,
          providers: [
            {
              inject: [ConfigService],
              provide: DATABASE_SERVICE_TOKEN, // Register a custom provider to override the default one
              useFactory: (configService: ConfigService) =>
                databaseServiceFactory(configService.get('database.url'))
            }
          ]
        }
      }
    },
    DatabaseModuleDefinition
  ))()
