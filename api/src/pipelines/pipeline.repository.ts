import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { CreatePipelineDto } from './dto/create-pipeline.dto'
import { UpdatePipelineDto } from './dto/update-pipeline.dto'
import { PipelineWithPipelineStepsModel } from './entities/pipeline.entity'

const PIPELINE_INCLUDE: Prisma.PipelineInclude = {
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
}

@Injectable()
export class PipelineRepository extends BaseRepository<
  PipelineWithPipelineStepsModel,
  CreatePipelineDto,
  UpdatePipelineDto,
  Prisma.PipelineInclude,
  Prisma.PipelineUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.pipeline, PIPELINE_INCLUDE)
  }

  async create(orgname: string, createPipelineDto: CreatePipelineDto) {
    const pipeline = await this.prisma.pipeline.create({
      data: {
        description: createPipelineDto.description,
        name: createPipelineDto.name,
        orgname
      },
      include: PIPELINE_INCLUDE
    })

    for (const pipelineStep of createPipelineDto.pipelineSteps) {
      await this.prisma.pipelineStep.create({
        data: {
          dependsOn: {
            connect: pipelineStep.dependsOn?.map((id) => ({
              id
            }))
          },
          id: pipelineStep.id,
          name: pipelineStep.name,
          pipelineId: pipeline.id,
          toolId: pipelineStep.toolId
        }
      })
    }

    return this.findOne(orgname, pipeline.id)
  }

  async createDefaultPipeline(orgname: string) {
    const pipeline = await this.prisma.pipeline.create({
      data: {
        description:
          'This is a default pipeline for indexing arbitrary documents. It extracts text from the document, creates an image from the text, summarizes the text, creates embeddings from the text, and converts the text to speech.',
        name: 'Default',
        orgname
      },
      include: PIPELINE_INCLUDE
    })
    const tools = await this.prisma.tool.findMany({
      where: {
        orgname
      }
    })

    // Create first step, this has no dependents
    const firstStep = await this.prisma.pipelineStep.create({
      data: {
        name: 'extract-text',
        pipelineId: pipeline.id,
        toolId: tools.find((t) => t.name == 'Extract Text').id
      }
    })
    const dependents = tools.filter((t) => t.name != 'Extract Text')

    for (const tool of dependents) {
      await this.prisma.pipelineStep.create({
        data: {
          dependsOn: {
            connect: {
              id: firstStep.id
            }
          },
          name: tool.toolBase,
          pipelineId: pipeline.id,
          toolId: tool.id
        }
      })
    }

    return this.findOne(orgname, pipeline.id)
  }

  async update(orgname: string, id: string, updatePipelineDto: UpdatePipelineDto) {
    const previousPipeline = await this.prisma.pipeline.findUnique({
      include: PIPELINE_INCLUDE,
      where: {
        id
      }
    })
    const pipelineStepsToDelete = previousPipeline.pipelineSteps.map((tool) => tool.id)

    await this.prisma.pipeline.update({
      data: {
        name: updatePipelineDto.name
      },
      include: PIPELINE_INCLUDE,
      where: {
        id
      }
    })

    await this.prisma.pipelineStep.deleteMany({
      where: {
        id: {
          in: pipelineStepsToDelete
        }
      }
    })

    for (const pipelineStep of updatePipelineDto.pipelineSteps) {
      await this.prisma.pipelineStep.create({
        data: {
          dependsOn: {
            connect: pipelineStep.dependsOn?.map((id) => ({
              id
            }))
          },
          id: pipelineStep.id,
          name: pipelineStep.name,
          pipelineId: id,
          toolId: pipelineStep.toolId
        }
      })
    }

    return this.findOne(orgname, id)
  }
}
