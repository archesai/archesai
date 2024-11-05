// run.repository.ts
import { Injectable } from "@nestjs/common";
import { Run, RunStatus } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { RunQueryDto } from "./dto/run-query.dto";

@Injectable()
export class RunRepository
  implements BaseRepository<Run, undefined, RunQueryDto, undefined>
{
  constructor(private readonly prisma: PrismaService) {}

  async createPipelineRun(
    orgname: string,
    pipelineId: string,
    runInputContentIds: string[]
  ) {
    const pipelineRun = await this.prisma.run.create({
      data: {
        name: "Pipeline Run",
        orgname,
        pipelineId,
        status: RunStatus.QUEUED,
        type: "PIPELINE_RUN",
      },
    });

    await this.prisma.runInputContent.createMany({
      data: runInputContentIds.map((contentId) => ({
        contentId,
        runId: pipelineRun.id,
      })),
    });

    // Step 3: Fetch the pipeline tools in order
    const pipelineTools = await this.prisma.pipelineTool.findMany({
      include: { dependsOn: true, tool: true },
      orderBy: {
        /* order as needed */
      },
      where: { pipelineId: "pipeline_id" },
    });

    // Step 4: Create child tool runs for each tool in the pipeline
    for (const pipelineTool of pipelineTools) {
      const toolRun = await this.prisma.run.create({
        data: {
          name: "Tool Run",
          orgname: orgname,
          parentRunId: pipelineRun.id,
          status: "QUEUED",
          toolId: pipelineTool.toolId,
          type: "TOOL_RUN",
        },
      });

      // Associate input contents for the tool run
      // This might depend on the output of the previous tool
      // For the first tool, it might use the pipeline's input contents
      let inputContentIds = [];

      if (pipelineTool.dependsOnId) {
        // Fetch output contents from the previous tool run
        const previousToolRun = await this.prisma.run.findFirst({
          include: {
            outputContents: true,
          },
          where: {
            parentRunId: pipelineRun.id,
            toolId: pipelineTool.dependsOn.toolId,
          },
        });
        inputContentIds = previousToolRun.outputContents.map(
          (outputContent) => outputContent.contentId
        );
      } else {
        // Use pipeline's input contents for the first tool
        const pipelineInputContents =
          await this.prisma.runInputContent.findMany({
            where: { runId: pipelineRun.id },
          });
        inputContentIds = pipelineInputContents.map(
          (inputContent) => inputContent.contentId
        );
      }

      // Associate input contents
      await this.prisma.runInputContent.createMany({
        data: inputContentIds.map((contentId) => ({
          contentId,
          runId: toolRun.id,
        })),
      });

      // After execution, associate output contents (this will happen after the tool run completes)
      // ...
    }

    return this.prisma.run.findUnique({
      where: { id: pipelineRun.id },
    });
  }

  async createToolRun(
    orgname: string,
    toolId: string,
    runInputContentIds: string[]
  ) {
    const run = await this.prisma.run.create({
      data: {
        name: "Tool Run",
        orgname,
        status: RunStatus.QUEUED,
        toolId,
        type: "TOOL_RUN",
      },
    });

    await this.prisma.runInputContent.createMany({
      data: runInputContentIds.map((contentId) => ({
        contentId,
        runId: run.id,
      })),
    });

    return this.prisma.run.findUnique({
      where: { id: run.id },
    });
  }

  async findAll(orgname: string, runQueryDto: RunQueryDto) {
    const count = await this.prisma.run.count({
      where: {
        createdAt: {
          gte: runQueryDto.startDate,
          lte: runQueryDto.endDate,
        },
        orgname,
      },
    });
    const results = await this.prisma.run.findMany({
      orderBy: {
        [runQueryDto.sortBy]: runQueryDto.sortDirection,
      },
      skip: runQueryDto.offset,
      take: runQueryDto.limit,
      where: {
        createdAt: {
          gte: runQueryDto.startDate,
          lte: runQueryDto.endDate,
        },
        orgname,
      },
    });
    return { count, results };
  }

  async findOne(orgname: string, id: string) {
    return this.prisma.run.findUniqueOrThrow({
      where: { id },
    });
  }

  async setCompletedAt(id: string, completedAt: Date) {
    return this.prisma.run.update({
      data: { completedAt },
      where: { id },
    });
  }

  async setProgress(id: string, progress: number) {
    return this.prisma.run.update({
      data: { progress },
      where: { id },
    });
  }

  async setRunError(id: string, error: string) {
    return this.prisma.run.update({
      data: { error },
      where: { id },
    });
  }

  async setStartedAt(id: string, startedAt: Date) {
    return this.prisma.run.update({
      data: { startedAt },
      where: { id },
    });
  }

  async updateStatus(id: string, status: RunStatus) {
    return this.prisma.run.update({
      data: { status },
      where: { id },
    });
  }
}
