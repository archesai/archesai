import { Injectable } from "@nestjs/common";
import { Pipeline, PipelineTool, Prisma, Tool } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";

@Injectable()
export class PipelineRepository extends BaseRepository<
  {
    pipelineTools: ({ tool: Tool } & PipelineTool)[];
  } & Pipeline,
  CreatePipelineDto,
  UpdatePipelineDto,
  Prisma.PipelineInclude,
  Prisma.PipelineSelect
> {
  constructor(private prisma: PrismaService) {
    super(prisma.pipeline, {
      pipelineTools: {
        include: {
          tool: true,
        },
      },
    });
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
      include: {
        pipelineTools: {
          include: {
            tool: true,
          },
        },
      },
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
}
