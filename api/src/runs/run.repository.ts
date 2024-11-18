import { Injectable } from "@nestjs/common";
import { Prisma, RunStatus, RunType } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { ContentEntity } from "../content/entities/content.entity";
import { PrismaService } from "../prisma/prisma.service";
import { CreateRunDto } from "./dto/create-run.dto";
import { RunModel } from "./entities/run.entity";

@Injectable()
export class RunRepository extends BaseRepository<
  RunModel,
  CreateRunDto,
  any,
  Prisma.RunInclude,
  Prisma.RunUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.run);
  }

  async createPipelineRun(
    orgname: string,
    pipelineId: string,
    createPipelineRunDto: CreateRunDto
  ) {
    const pipeline = await this.prisma.pipeline.findUniqueOrThrow({
      include: {
        pipelineSteps: true,
      },
      where: { id: pipelineId },
    });
    const pipelineRun = await this.prisma.run.create({
      data: {
        name: "Pipeline Run",
        orgname,
        pipelineId,
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
              status: RunStatus.QUEUED,
            })),
          },
        },
      },
    });

    for (const pipelineStep of pipeline.pipelineSteps) {
      await this.prisma.run.update({
        data: {
          inputs: {
            connect: createPipelineRunDto.contentIds.map((contentId) => ({
              id: contentId,
            })),
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

    return this.prisma.run.findUnique({
      where: { id: pipelineRun.id },
    });
  }

  async setOutputContent(toolRunId: string, contents: ContentEntity[]) {
    await this.prisma.run.update({
      data: {
        inputs: {
          connect: contents.map((content) => ({ id: content.id })),
        },
      },
      where: { id: toolRunId },
    });
    return this.prisma.run.findUnique({
      where: { id: toolRunId },
    });
  }
}
