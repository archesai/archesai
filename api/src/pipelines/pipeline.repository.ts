import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'

const PIPELINE_INCLUDE = {
  pipelineSteps: {
    include: {
      dependents: {
        select: {
          id: true,
          name: true
        }
      },
      dependsOn: {
        select: {
          id: true,
          name: true
        }
      },
      tool: true
    }
  }
} as const

@Injectable()
export class PipelineRepository extends BaseRepository<
  Prisma.PipelineDelegate,
  typeof PIPELINE_INCLUDE
> {
  constructor(private prisma: PrismaService) {
    super(prisma.pipeline, PIPELINE_INCLUDE)
  }
}

@Injectable()
export class PipelineStepRepository extends BaseRepository<Prisma.PipelineStepDelegate> {
  constructor(private prisma: PrismaService) {
    super(prisma.pipelineStep)
  }
}
