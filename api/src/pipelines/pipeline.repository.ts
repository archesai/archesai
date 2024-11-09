import { Injectable } from "@nestjs/common";
import { Prisma, RunStatus } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { CreatePipelineRunDto } from "./dto/create-pipeline-run.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import { PipelineWithPipelineStepsModel } from "./entities/pipeline.entity";
import { PipelineRunEntity } from "./entities/pipeline-run.entity";

const PIPELINE_INCLUDE = {
  pipelineSteps: {
    include: {
      tool: true,
    },
  },
};

@Injectable()
export class PipelineRepository extends BaseRepository<
  PipelineWithPipelineStepsModel,
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
        pipelineSteps: {
          createMany: {
            data: createPipelineDto.pipelineSteps.map((tool) => {
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

  async createPipelineRun(
    orgname: string,
    pipelineId: string,
    createPipelineRunDto: CreatePipelineRunDto
  ) {
    const pipeline = await this.findOne(orgname, pipelineId);
    const pipelineRun = await this.prisma.pipelineRun.create({
      data: {
        name: "Pipeline Run",
        orgname,
        pipelineId,
        status: RunStatus.QUEUED,
        transformations: {
          createMany: {
            data: pipeline.pipelineSteps.map((pipelineStep) => ({
              createdAt: new Date(),
              name: new Date().toISOString(),
              pipelineStepId: pipelineStep.id,
              status: RunStatus.QUEUED,
            })),
          },
        },
      },
    });

    for (const pipelineStep of pipeline.pipelineSteps) {
      await this.prisma.transformation.update({
        data: {
          inputs: {
            connect: createPipelineRunDto.runInputContentIds.map(
              (contentId) => ({
                id: contentId,
              })
            ),
          },
        },
        where: {
          pipelineRunId_pipelineStepId: {
            pipelineRunId: pipelineRun.id,
            pipelineStepId: pipelineStep.id,
          },
          pipelineStep: {
            dependsOnId: null,
          },
        },
      });
    }

    return new PipelineRunEntity(
      await this.prisma.pipelineRun.findUnique({
        where: { id: pipelineRun.id },
      })
    );
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
    const pipelineStepsToDelete = previousPipeline.pipelineSteps.map(
      (tool) => tool.id
    );

    return this.prisma.pipeline.update({
      data: {
        name: updatePipelineDto.name,
        ...(updatePipelineDto.pipelineSteps
          ? {
              pipelineSteps: {
                createMany: {
                  data: updatePipelineDto.pipelineSteps.map((tool) => {
                    return {
                      dependsOnId: tool.dependsOnId,
                      toolId: tool.toolId,
                    };
                  }),
                },
                deleteMany: {
                  id: {
                    in: pipelineStepsToDelete,
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
