import { Injectable, Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import {
  PipelineEntity,
  PipelineWithPipelineStepsModel,
} from "./entities/pipeline.entity";
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
    private websocketsService: WebsocketsService
  ) {
    super(pipelineRepository);
  }

  async createDefaultPipeline(orgname: string) {
    return this.toEntity(
      await this.pipelineRepository.createDefaultPipeline(orgname)
    );
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "pipelines"],
    });
  }

  protected toEntity(model: PipelineWithPipelineStepsModel): PipelineEntity {
    return new PipelineEntity(model);
  }
}
