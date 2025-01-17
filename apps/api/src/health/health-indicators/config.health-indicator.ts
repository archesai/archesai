import { ConfigService } from '@/src/config/config.service'
import { RunStatusEnum } from '@/src/runs/entities/run.entity'
import { Injectable } from '@nestjs/common'
import {
  HealthIndicator,
  HealthIndicatorResult,
  HealthCheckError
} from '@nestjs/terminus'

@Injectable()
export class ConfigHealthIndicator extends HealthIndicator {
  constructor(private configService: ConfigService) {
    super()
  }
  async isHealthy(key: string): Promise<HealthIndicatorResult> {
    const health = this.configService.getHealth()
    const isHealthy = health.status === RunStatusEnum.COMPLETE
    const status = this.getStatus(key, isHealthy, { health })

    if (isHealthy) {
      return status
    }
    throw new HealthCheckError('config health check failed', status)
  }
}
