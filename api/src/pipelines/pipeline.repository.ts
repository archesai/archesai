import { Injectable } from "@nestjs/common";
import { Pipeline, Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { PipelineQueryDto } from "./dto/pipeline-query.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";

@Injectable()
export class PipelineRepository
  implements
    BaseRepository<
      Pipeline,
      CreatePipelineDto,
      PipelineQueryDto,
      UpdatePipelineDto
    >
{
  constructor(private prisma: PrismaService) {}

  async create(orgname: string, createPipelineDto: CreatePipelineDto) {
    return this.prisma.pipeline.create({
      data: {
        name: createPipelineDto.name,
        organization: {
          connect: {
            orgname,
          },
        },
        pipelineTools: {
          createMany: {
            data: createPipelineDto.pipelineTools.map((tool) => {
              return {
                dependsOnId: tool.dependsOnId,
                toolId: tool.toolId,
              };
            }),
          },
        },
      },
      include: {
        pipelineTools: {
          include: {
            tool: true,
          },
        },
      },
    });
  }

  async findAll(orgname: string, pipelineQueryDto: PipelineQueryDto) {
    const whereConditions = {
      createdAt: {
        gte: pipelineQueryDto.startDate,
        lte: pipelineQueryDto.endDate,
      },
      orgname,
    };
    if (pipelineQueryDto.filters) {
      pipelineQueryDto.filters.forEach((filter) => {
        whereConditions[filter.field] = { [filter.operator]: filter.value };
      });
    }

    const count = await this.prisma.pipeline.count({
      where: whereConditions,
    });
    const pipelines = await this.prisma.pipeline.findMany({
      include: {
        pipelineTools: {
          include: {
            tool: true,
          },
        },
      },
      orderBy: {
        [pipelineQueryDto.sortBy]: pipelineQueryDto.sortDirection,
      },
      skip: pipelineQueryDto.offset,
      take: pipelineQueryDto.limit,
      where: whereConditions,
    });
    return { count, results: pipelines };
  }

  async findOne(id: string) {
    return this.prisma.pipeline.findUniqueOrThrow({
      include: {
        pipelineTools: {
          include: {
            tool: true,
          },
        },
      },
      where: { id },
    });
  }

  async remove(orgname: string, id: string) {
    await this.prisma.pipeline.delete({
      where: { id },
    });
  }

  async update(
    orgname: string,
    id: string,
    updatePipelineDto: UpdatePipelineDto
  ) {
    const previousPipeline = await this.prisma.pipeline.findUnique({
      select: {
        pipelineTools: {
          include: {
            tool: true,
          },
        },
      },
      where: {
        id,
      },
    });
    const pipelineToolsToDelete = previousPipeline.pipelineTools.map(
      (tool) => tool.id
    );
    return this.prisma.pipeline.update({
      data: {
        name: updatePipelineDto.name,

        ...(updatePipelineDto.pipelineTools
          ? {
              pipelineTools: {
                createMany: {
                  data: updatePipelineDto.pipelineTools.map((tool) => {
                    return {
                      dependsOnId: tool.dependsOnId,
                      toolId: tool.toolId,
                    };
                  }),
                },
                deleteMany: {
                  id: {
                    in: pipelineToolsToDelete,
                  },
                },
              },
            }
          : {}),

        //
      },
      include: {
        pipelineTools: {
          include: {
            tool: true,
          },
        },
      },
      where: {
        id,
      },
    });
  }

  async updateRaw(
    orgname: string,
    id: string,
    raw: Prisma.PipelineUpdateInput
  ) {
    return this.prisma.pipeline.update({
      data: raw,
      include: {
        pipelineTools: {
          include: {
            tool: true,
          },
        },
      },
      where: {
        id,
      },
    });
  }
}
