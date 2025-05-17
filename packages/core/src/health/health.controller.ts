// import {
//   DiskHealthIndicator,
//   HealthCheck,
//   HealthCheckService,
//   HttpHealthIndicator,
//   MemoryHealthIndicator
// } from '/_terminus'

// import { DatabaseService } from '@archesai/core/drizzle/drizzle.service'
// import { ConfigHealthIndicator } from '#config/config.health-indicator'

// import { LlmHealthIndicator } from '#health/health-indicators/llm.health-indicator'
// import { StorageHealthIndicator } from '#health/health-indicators/storage.health-indicator'

export class HealthController {
  // private readonly config: ConfigHealthIndicator
  // private readonly disk: DiskHealthIndicator
  // private readonly health: HealthCheckService
  // private readonly http: HttpHealthIndicator
  // private readonly memory: MemoryHealthIndicator
  // constructor(
  //   health: HealthCheckService,
  //   disk: DiskHealthIndicator,
  //   config: ConfigHealthIndicator,
  //   http: HttpHealthIndicator,
  //   memory: MemoryHealthIndicator
  // ) {
  //   this.health = health
  //   this.disk = disk
  //   this.config = config
  //   this.http = http
  //   this.memory = memory
  // }
  // @ApiOkResponse({
  //   type: HealthCheck
  // })
  // @ApiOperation({
  //   description: 'This endpoint will check the health of the application',
  //   summary: 'Health Check'
  // })
  // @Get()
  // @HealthCheck()
  // check() {
  //   return this.health.check([
  //     () =>
  //       this.disk.checkStorage('storage', { path: '/', thresholdPercent: 0.5 }),
  //     () => this.config.isHealthy('config'),
  //     () => this.memory.checkHeap('memory_heap', 150 * 1024 * 1024),
  //     () => this.memory.checkRSS('memory_rss', 150 * 1024 * 1024)
  //     // () => this.llm.isHealthy('llm'),
  //     // () => this.storage.isHealthy('storage')
  //   ])
  // }
}
