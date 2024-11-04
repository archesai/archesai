import { InjectFlowProducer } from "@nestjs/bullmq";
import { Injectable, Logger } from "@nestjs/common";
import { Prisma } from "@prisma/client";
import { FlowProducer } from "bullmq";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { ContentEntity } from "../content/entities/content.entity";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { PipelineQueryDto } from "./dto/pipeline-query.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import { PipelineEntity } from "./entities/pipeline.entity";
import { PipelineRepository } from "./pipeline.repository";

@Injectable()
export class PipelinesService
  implements
    BaseService<
      PipelineEntity,
      CreatePipelineDto,
      PipelineQueryDto,
      UpdatePipelineDto
    >
{
  private logger = new Logger(PipelinesService.name);

  constructor(
    private pipelineRepository: PipelineRepository,
    private websocketsService: WebsocketsService,
    @InjectFlowProducer("flow") private readonly flowProducer: FlowProducer
  ) {}

  async create(
    orgname: string,
    createPipelineDto: CreatePipelineDto
  ): Promise<PipelineEntity> {
    const pipeline = await this.pipelineRepository.create(
      orgname,
      createPipelineDto
    );
    return new PipelineEntity(pipeline);
  }

  async findAll(orgname: string, pipelineQueryDto: PipelineQueryDto) {
    const { count, results } = await this.pipelineRepository.findAll(
      orgname,
      pipelineQueryDto
    );
    const pipelineEntities = results.map(
      (pipeline) => new PipelineEntity(pipeline)
    );
    return new PaginatedDto<PipelineEntity>({
      metadata: {
        limit: pipelineQueryDto.limit,
        offset: pipelineQueryDto.offset,
        totalResults: count,
      },
      results: pipelineEntities,
    });
  }

  async findOne(id: string) {
    const pipeline = await this.pipelineRepository.findOne(id);
    return new PipelineEntity(pipeline);
  }

  async remove(orgname: string, pipelineId: string): Promise<void> {
    await this.pipelineRepository.remove(orgname, pipelineId);
    this.websocketsService.socket.to(orgname).emit("update");
  }

  async runPipeline(content: ContentEntity) {
    await this.flowProducer.add({
      data: {
        content: content,
        toolId: "extract-text",
      },
      name: "extract-text",
      queueName: "tool",
    });
  }

  async update(
    orgname: string,
    id: string,
    updatePipelineDto: UpdatePipelineDto
  ) {
    const pipeline = await this.pipelineRepository.update(
      orgname,
      id,
      updatePipelineDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return new PipelineEntity(pipeline);
  }

  async updateRaw(
    orgname: string,
    id: string,
    raw: Prisma.PipelineUpdateInput
  ) {
    const pipeline = await this.pipelineRepository.updateRaw(orgname, id, raw);
    this.websocketsService.socket.to(orgname).emit("update");
    return new PipelineEntity(pipeline);
  }
}
