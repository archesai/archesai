import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, createModule } from '@archesai/core'

import { GoogleCloudStorageService } from '#storage/services/gcp.service'
import { LocalStorageService } from '#storage/services/local.service'
import { S3StorageProvider } from '#storage/services/s3.service'
import { StorageService } from '#storage/storage.service'

export const StorageModuleDefinition: ModuleMetadata = {
  exports: [StorageService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: StorageService,
      useFactory: (configService: ConfigService): StorageService => {
        const storageType = configService.get('storage.type')
        switch (storageType) {
          case 'google-cloud':
            return new GoogleCloudStorageService()
          case 'local':
            return new LocalStorageService()
          case 'minio':
            return new S3StorageProvider(configService)
          default:
            return new GoogleCloudStorageService()
        }
      }
    }
  ]
}

export const StorageModule = (() =>
  createModule(class StorageModule {}, StorageModuleDefinition))()
