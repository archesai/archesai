import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import { PipelineWithPipelineToolsModel } from "./entities/pipeline.entity";

const PIPELINE_INCLUDE = {
  pipelineTools: {
    include: {
      tool: true,
    },
  },
};

@Injectable()
export class PipelineRepository extends BaseRepository<
  PipelineWithPipelineToolsModel,
  CreatePipelineDto,
  UpdatePipelineDto,
  Prisma.PipelineInclude,
  Prisma.PipelineSelect,
  Prisma.PipelineUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.pipeline, PIPELINE_INCLUDE);
  }

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
      include: PIPELINE_INCLUDE,
    });
  }

  async update(
    orgname: string,
    id: string,
    updatePipelineDto: UpdatePipelineDto
  ) {
    const previousPipeline = await this.prisma.pipeline.findUnique({
      include: PIPELINE_INCLUDE,
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
      },
      include: PIPELINE_INCLUDE,
      where: {
        id,
      },
    });
  }
}
