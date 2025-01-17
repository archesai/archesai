import { RunStatusEnum } from '@/src/runs/entities/run.entity'
import {
  IStorageService,
  STORAGE_SERVICE
} from '@/src/storage/interfaces/storage-provider.interface'
import { Inject, Injectable } from '@nestjs/common'
import {
  HealthIndicator,
  HealthIndicatorResult,
  HealthCheckError
} from '@nestjs/terminus'

@Injectable()
export class StorageHealthIndicator extends HealthIndicator {
  constructor(
    @Inject(STORAGE_SERVICE)
    private storageService: IStorageService
  ) {
    super()
  }
  async isHealthy(key: string): Promise<HealthIndicatorResult> {
    const health = this.storageService.getHealth()
    const isHealthy = health.status === RunStatusEnum.COMPLETE
    const status = this.getStatus(key, isHealthy, { health })

    if (isHealthy) {
      return status
    }
    throw new HealthCheckError('storage health check failed', status)
  }
}
