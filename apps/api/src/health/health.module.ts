import { Module } from '@nestjs/common'
import { TerminusModule } from '@nestjs/terminus'
import { HttpModule } from '@nestjs/axios'
import { HealthController } from './health.controller'
import { PrismaModule } from '../prisma/prisma.module'
import { ConfigHealthIndicator } from '@/src/health/health-indicators/config.health-indicator'
import { StorageHealthIndicator } from '@/src/health/health-indicators/storage.health-indicator'
import { LlmModule } from '@/src/llm/llm.module'
import { StorageModule } from '@/src/storage/storage.module'
import { LlmHealthIndicator } from '@/src/health/health-indicators/llm.health-indicator'

@Module({
  imports: [
    TerminusModule,
    HttpModule,
    PrismaModule,
    LlmModule,
    StorageModule.forRoot()
  ],
  controllers: [HealthController],
  providers: [ConfigHealthIndicator, StorageHealthIndicator, LlmHealthIndicator]
})
export class HealthModule {}
