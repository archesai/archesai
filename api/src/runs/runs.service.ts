import { InjectFlowProducer } from "@nestjs/bullmq";
import { BadRequestException, Injectable, Logger } from "@nestjs/common";
import { RunStatus, RunType } from "@prisma/client";
import { FlowProducer } from "bullmq";

import { BaseService } from "../common/base.service";
import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { PipelinesService } from "../pipelines/pipelines.service";
import { ToolsService } from "../tools/tools.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateRunDto } from "./dto/create-run.dto";
import { RunEntity, RunModel } from "./entities/run.entity";
import { RunRepository } from "./run.repository";

@Injectable()
export class RunsService extends BaseService<
  RunEntity,
  CreateRunDto,
  any,
  RunRepository,
  RunModel
> {
  private logger = new Logger(RunsService.name);

  constructor(
    private runRepository: RunRepository,
    private websocketsService: WebsocketsService,
    private pipelinesService: PipelinesService,
    private toolsService: ToolsService,
    @InjectFlowProducer("flow") private readonly flowProducer: FlowProducer,
    private contentService: ContentService
  ) {
    super(runRepository);
  }

  async createPipelineRun(
    orgname: string,
    pipelineId: string,
    createPipelineRunDto: CreateRunDto
  ) {
    // Ensure run content
    const runContent = await this.ensureRunContent(
      orgname,
      pipelineId,
      createPipelineRunDto
    );

    // Create the run in the database
    const run = await this.runRepository.createPipelineRun(
      orgname,
      pipelineId,
      {
        contentIds: runContent.map((content) => content.id),
        pipelineId,
        runType: RunType.PIPELINE_RUN,
      }
    );
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "pipelines"],
    });

    // Add the run to the queue
    await this.flowProducer.add({
      data: {
        // content: content,
        toolId: "extract-text",
      },
      name: "extract-text",
      queueName: "tool",
    });
    return this.toEntity(run);
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "tool-runs"],
    });
  }

  async ensureRunContent(
    orgname: string,
    pipelineId: string,
    createPipelineRunDto: CreateRunDto
  ) {
    const pipeline = await this.pipelinesService.findOne(orgname, pipelineId);

    const runContent: ContentEntity[] = [];
    if (!!createPipelineRunDto.contentIds?.length) {
      for (const contentId of createPipelineRunDto.contentIds) {
        runContent.push(await this.contentService.findOne(orgname, contentId));
      }
    }
    if (createPipelineRunDto.text) {
      runContent.push(
        await this.contentService.create(orgname, {
          name: "Pipeline Input Text - " + pipeline.id,
          text: createPipelineRunDto.text,
        })
      );
    }
    if (createPipelineRunDto.url) {
      runContent.push(
        await this.contentService.create(orgname, {
          name: "Tool Input URL - " + pipeline.id,
          url: createPipelineRunDto.url,
        })
      );
    }
    if (!runContent.length) {
      throw new BadRequestException("No input content provided");
    }

    return runContent;
  }

  async setOutputContent(toolRunId: string, content: ContentEntity[]) {
    return this.toEntity(
      await this.runRepository.setOutputContent(toolRunId, content)
    );
  }

  async setProgress(id: string, progress: number) {
    return this.toEntity(
      await this.runRepository.updateRaw(null, id, {
        progress,
      })
    );
  }

  async setRunError(id: string, error: string) {
    return this.toEntity(
      await this.runRepository.updateRaw(null, id, {
        error,
      })
    );
  }

  async setStatus(id: string, status: RunStatus) {
    switch (status) {
      case "COMPLETE":
        await this.runRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        await this.runRepository.updateRaw(null, id, {
          progress: 1,
        });
        break;
      case "ERROR":
        await this.runRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        break;
      case "PROCESSING":
        await this.runRepository.updateRaw(null, id, {
          startedAt: new Date(),
        });
        break;
    }
    const run = await this.runRepository.updateRaw(null, id, {
      status,
    });
    return this.toEntity(run);
  }

  protected toEntity(model: RunModel): RunEntity {
    return new RunEntity(model);
  }
}
