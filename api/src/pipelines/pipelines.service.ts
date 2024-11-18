import { InjectFlowProducer } from "@nestjs/bullmq";
import { BadRequestException, Injectable, Logger } from "@nestjs/common";
import { FlowProducer } from "bullmq";

import { BaseService } from "../common/base.service";
import { CreateRunDto } from "../common/dto/create-run.dto";
import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import {
  PipelineEntity,
  PipelineWithPipelineStepsModel,
} from "./entities/pipeline.entity";
import { PipelineRunEntity } from "./entities/pipeline-run.entity";
import { PipelineRepository } from "./pipeline.repository";

@Injectable()
export class PipelinesService extends BaseService<
  PipelineEntity,
  CreatePipelineDto,
  UpdatePipelineDto,
  PipelineRepository,
  PipelineWithPipelineStepsModel
> {
  private logger = new Logger(PipelinesService.name);

  constructor(
    private pipelineRepository: PipelineRepository,
    @InjectFlowProducer("flow") private readonly flowProducer: FlowProducer,
    private websocketsService: WebsocketsService,
    private contentService: ContentService
  ) {
    super(pipelineRepository);
  }

  async createDefaultPipeline(orgname: string) {
    return this.toEntity(
      await this.pipelineRepository.createDefaultPipeline(orgname)
    );
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
    const run = await this.pipelineRepository.createPipelineRun(
      orgname,
      pipelineId,
      {
        contentIds: runContent.map((content) => content.id),
      }
    );
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "piplines"],
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
    return new PipelineRunEntity(run);
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "pipelines"],
    });
  }

  async ensureRunContent(
    orgname: string,
    pipelineId: string,
    createPipelineRunDto: CreateRunDto
  ) {
    const pipeline = await this.pipelineRepository.findOne(orgname, pipelineId);

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

  protected toEntity(model: PipelineWithPipelineStepsModel): PipelineEntity {
    return new PipelineEntity(model);
  }
}
