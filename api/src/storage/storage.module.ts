import { DynamicModule, Module } from '@nestjs/common'

import { StorageController } from './storage.controller'
import { GoogleCloudStorageService } from './storage.gcp.service'
import { LocalStorageService } from './storage.local.service'
import { S3StorageProvider } from './storage.s3.service'
import { STORAGE_SERVICE, StorageService } from './storage.service'
import { ArchesConfigService } from '../config/config.service'

@Module({})
export class StorageModule {
  static forRoot(): DynamicModule {
    return {
      controllers: [StorageController],
      exports: [STORAGE_SERVICE],
      module: StorageModule,
      providers: [
        {
          inject: [ArchesConfigService],
          provide: STORAGE_SERVICE,
          useFactory: (configService: ArchesConfigService): StorageService => {
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
  }
}
