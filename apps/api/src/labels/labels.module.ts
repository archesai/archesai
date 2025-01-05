import { Module } from '@nestjs/common'

import { PrismaModule } from '../prisma/prisma.module'
import { LabelRepository } from './label.repository'
import { LabelsController } from './labels.controller'
import { LabelsService } from './labels.service'

@Module({
  controllers: [LabelsController],
  exports: [LabelsService],
  imports: [PrismaModule],
  providers: [LabelsService, LabelRepository]
})
export class LabelsModule {}
