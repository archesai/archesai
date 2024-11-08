import { Injectable, Logger } from "@nestjs/common";
import { Pipeline, PipelineTool, Tool } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import { PipelineEntity } from "./entities/pipeline.entity";
import { PipelineRepository } from "./pipeline.repository";

@Injectable()
export class PipelinesService extends BaseService<
  PipelineEntity,
  CreatePipelineDto,
  UpdatePipelineDto,
  PipelineRepository,
  {
    pipelineTools: ({ tool: Tool } & PipelineTool)[];
  } & Pipeline
> {
  private logger = new Logger(PipelinesService.name);

  constructor(private pipelineRepository: PipelineRepository) {
    super(pipelineRepository);
  }

  protected toEntity(
    model: {
      pipelineTools: ({ tool: Tool } & PipelineTool)[];
    } & Pipeline
  ): PipelineEntity {
    return new PipelineEntity(model);
  }
}
