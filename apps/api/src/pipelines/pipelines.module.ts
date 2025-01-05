import { Module } from '@nestjs/common'

import { PrismaModule } from '../prisma/prisma.module'
import {
  PipelineRepository,
  PipelineStepRepository
} from './pipeline.repository'
import { PipelinesController } from './pipelines.controller'
import { PipelinesService } from './pipelines.service'
import { ToolsModule } from '../tools/tools.module'

@Module({
  controllers: [PipelinesController],
  exports: [PipelinesService],
  imports: [PrismaModule, ToolsModule],
  providers: [PipelinesService, PipelineRepository, PipelineStepRepository]
})
export class PipelinesModule {}
