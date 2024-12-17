import { Module } from '@nestjs/common'
import { TerminusModule } from '@nestjs/terminus'
import { HttpModule } from '@nestjs/axios'
import { HealthController } from './health.controller'
import { PrismaModule } from '../prisma/prisma.module'

@Module({
  imports: [TerminusModule, HttpModule, PrismaModule],
  controllers: [HealthController]
})
export class HealthModule {}
