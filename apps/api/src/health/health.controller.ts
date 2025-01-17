import { Controller, Get } from '@nestjs/common'
import {
  HealthCheckService,
  HttpHealthIndicator,
  HealthCheck,
  DiskHealthIndicator,
  PrismaHealthIndicator
} from '@nestjs/terminus'
import { PrismaService } from '@/src/prisma/prisma.service'
import { ConfigHealthIndicator } from '@/src/health/health-indicators/config.health-indicator'
import { LlmHealthIndicator } from '@/src/health/health-indicators/llm.health-indicator'
import { StorageHealthIndicator } from '@/src/health/health-indicators/storage.health-indicator'

@Controller('health')
export class HealthController {
  constructor(
    private readonly health: HealthCheckService,
    private readonly http: HttpHealthIndicator,
    private readonly disk: DiskHealthIndicator,
    private readonly prismaHealth: PrismaHealthIndicator,
    private readonly prisma: PrismaService,
    private readonly configHealthIndicator: ConfigHealthIndicator,
    private readonly llmHealthIndicator: LlmHealthIndicator,
    private readonly storageIndicator: StorageHealthIndicator
  ) {}

  @Get()
  @HealthCheck()
  check() {
    return this.health.check([
      () => this.http.pingCheck('nestjs-docs', 'https://docs.nestjs.com'),
      () =>
        this.disk.checkStorage('storage', { path: '/', thresholdPercent: 0.5 }),
      () => this.prismaHealth.pingCheck('db', this.prisma),
      () => this.configHealthIndicator.isHealthy('config'),
      () => this.llmHealthIndicator.isHealthy('llm'),
      () => this.storageIndicator.isHealthy('storage')
    ])
  }
}
