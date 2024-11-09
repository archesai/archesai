import { Injectable } from "@nestjs/common";
import { Prisma, RunStatus } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { CreatePipelineRunDto } from "./dto/create-pipeline-run.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import { PipelineWithPipelineStepsModel } from "./entities/pipeline.entity";

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
    const pipelineRun = await this.prisma.pipelineRun.create({
      data: {
        name: "Pipeline Run",
        orgname,
        pipelineId,
        status: RunStatus.QUEUED,
      },
    });

    await this.prisma.runContent.createMany({
      data: createPipelineRunDto.runInputContentIds.map((contentId) => ({
        contentId,
        pipelineRunId: pipelineRun.id,
        role: "INPUT",
      })),
    });

    // Step 3: Fetch the pipeline tools in order
    const pipelineSteps = await this.prisma.pipelineStep.findMany({
      include: { dependsOn: true, tool: true },
      orderBy: {
        /* order as needed */
      },
      where: { pipelineId },
    });

    // Step 4: Create child tool runs for each tool in the pipeline
    for (const pipelineStep of pipelineSteps) {
      const createdAt = new Date();
      const pipelineStepRun = await this.prisma.transformation.create({
        data: {
          createdAt,
          name: createdAt.toISOString(),
          pipelineRunId: pipelineRun.id,
          pipelineStepId: pipelineStep.id,
          status: "QUEUED",
        },
      });
      if (!pipelineStep.dependsOn) {
        // If there are no dependencies, the tool can be run immediately
        await this.prisma.runContent.createMany({
          data: createPipelineRunDto.runInputContentIds.map((contentId) => ({
            contentId,
            pipelineRunId: pipelineRun.id,
            piplineStepRunId: pipelineStepRun.id,
            role: "INPUT",
          })),
        });
      }
    }

    return this.prisma.pipelineRun.findUnique({
      where: { id: pipelineRun.id },
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
