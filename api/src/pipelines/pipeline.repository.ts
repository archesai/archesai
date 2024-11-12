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

  async createDefaultPipeline(orgname: string) {
    const pipeline = await this.prisma.pipeline.create({
      data: {
        description:
          "This is a default pipeline for indexing arbitrary documents. It extracts text from the document, creates an image from the text, summarizes the text, creates embeddings from the text, and converts the text to speech.",
        name: "Default",
        orgname,
      },
      include: PIPELINE_INCLUDE,
    });
    const tools = await this.prisma.tool.findMany({
      where: {
        orgname,
      },
    });

    // Create first step, this has no dependents
    const firstStep = await this.prisma.pipelineStep.create({
      data: {
        pipelineId: pipeline.id,
        toolId: tools.find((t) => t.name == "Extract Text").id,
      },
    });
    const dependents = tools.filter((t) => t.name != "Extract Text");

    for (const tool of dependents) {
      await this.prisma.pipelineStep.create({
        data: {
          dependsOn: {
            connect: {
              id: firstStep.id,
            },
          },
          pipelineId: pipeline.id,
          toolId: tool.id,
        },
      });
    }

    return this.findOne(orgname, pipeline.id);
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
        toolRuns: {
          createMany: {
            data: pipeline.pipelineSteps.map((pipelineStep) => ({
              createdAt: new Date(),
              name: new Date().toISOString(),
              orgname,
              pipelineStepId: pipelineStep.id,
              status: RunStatus.QUEUED,
            })),
          },
        },
      },
    });

    for (const pipelineStep of pipeline.pipelineSteps) {
      await this.prisma.toolRun.update({
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
            dependsOn: {
              none: {},
            },
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
