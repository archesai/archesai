import { DynamicModule, Module } from '@nestjs/common'

import { StorageController } from './storage.controller'
import { GoogleCloudStorageService } from './services/gcp.service'
import { LocalStorageService } from './services/local.service'
import { S3StorageProvider } from './services/s3.service'
import {
  IStorageService,
  STORAGE_SERVICE
} from './interfaces/storage-provider.interface'

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
          useFactory: (configService: ArchesConfigService): IStorageService => {
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
