import { Injectable } from '@nestjs/common'
import { Prisma, RunStatus, RunType } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { ContentEntity } from '../content/entities/content.entity'
import { PrismaService } from '../prisma/prisma.service'
import { CreateRunDto } from './dto/create-run.dto'
import { RunModel } from './entities/run.entity'

const RUN_INCLUDE = {
  inputs: {
    select: {
      id: true,
      name: true
    }
  },
  outputs: {
    select: {
      id: true,
      name: true
    }
  }
  // pipeline: true,
  // pipelineRun: true,
  // pipelineStep: true,
  // tool: true,
  // toolRuns: true,
}

@Injectable()
export class RunRepository extends BaseRepository<
  RunModel,
  CreateRunDto,
  any,
  Prisma.RunInclude,
  Prisma.RunUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.run, RUN_INCLUDE)
  }

  async createPipelineRun(orgname: string, createRunDto: CreateRunDto) {
    const pipeline = await this.prisma.pipeline.findUniqueOrThrow({
      include: {
        pipelineSteps: true
      },
      where: { id: createRunDto.pipelineId }
    })
    const pipelineRun = await this.prisma.run.create({
      data: {
        name: 'Pipeline Run',
        orgname,
        pipelineId: createRunDto.pipelineId,
        runType: RunType.PIPELINE_RUN,
        status: RunStatus.QUEUED,
        toolRuns: {
          createMany: {
            data: pipeline.pipelineSteps.map((pipelineStep) => ({
              createdAt: new Date(),
              name: new Date().toISOString(),
              orgname,
              pipelineId: pipeline.id,
              pipelineStepId: pipelineStep.id,
              runType: RunType.TOOL_RUN,
              status: RunStatus.QUEUED
            }))
          }
        }
      }
    })

    for (const pipelineStep of pipeline.pipelineSteps) {
      await this.prisma.run.update({
        data: {
          inputs: {
            connect: createRunDto.contentIds.map((contentId) => ({
              id: contentId
            }))
          }
        },
        where: {
          pipelineRunId_pipelineStepId: {
            pipelineRunId: pipelineRun.id,
            pipelineStepId: pipelineStep.id
          },
          pipelineStep: {
            dependsOn: {
              none: {}
            }
          }
        }
      })
    }

    return this.prisma.run.findUnique({
      include: RUN_INCLUDE,
      where: { id: pipelineRun.id }
    })
  }

  async createToolRun(orgname: string, createRunDto: CreateRunDto) {
    return this.prisma.run.create({
      data: {
        name: 'Tool Run',
        orgname,
        runType: RunType.TOOL_RUN,
        status: RunStatus.QUEUED,
        toolId: createRunDto.toolId
      },
      include: RUN_INCLUDE
    })
  }

  async setInputsOrOutputs(
    runId: string,
    type: 'inputs' | 'outputs',
    contents: ContentEntity[]
  ) {
    return this.prisma.run.update({
      data: {
        [type]: {
          connect: contents.map((content) => ({ id: content.id }))
        }
      },
      include: RUN_INCLUDE,
      where: { id: runId }
    })
  }
}
