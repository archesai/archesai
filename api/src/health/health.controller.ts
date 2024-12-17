import { Controller, Get } from '@nestjs/common'
import {
  HealthCheckService,
  HttpHealthIndicator,
  HealthCheck,
  DiskHealthIndicator,
  PrismaHealthIndicator
} from '@nestjs/terminus'
import { PrismaService } from '@/src/prisma/prisma.service'
import { IsPublic } from '@/src/auth/decorators/is-public.decorator'

@IsPublic()
@Controller('health')
export class HealthController {
  constructor(
    private health: HealthCheckService,
    private http: HttpHealthIndicator,
    private readonly disk: DiskHealthIndicator,
    private readonly prismaHealth: PrismaHealthIndicator,
    private prisma: PrismaService
  ) {}

  @Get()
  @HealthCheck()
  check() {
    return this.health.check([
      () => this.http.pingCheck('nestjs-docs', 'https://docs.nestjs.com'),
      () =>
        this.disk.checkStorage('storage', { path: '/', thresholdPercent: 0.5 }),
      () => this.prismaHealth.pingCheck('db', this.prisma)
    ])
  }
}
