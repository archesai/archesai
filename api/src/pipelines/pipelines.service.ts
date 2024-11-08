import { Injectable, Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import {
  PipelineEntity,
  PipelineWithPipelineToolsModel,
} from "./entities/pipeline.entity";
import { PipelineRepository } from "./pipeline.repository";

@Injectable()
export class PipelinesService extends BaseService<
  PipelineEntity,
  CreatePipelineDto,
  UpdatePipelineDto,
  PipelineRepository,
  PipelineWithPipelineToolsModel
> {
  private logger = new Logger(PipelinesService.name);

  constructor(private pipelineRepository: PipelineRepository) {
    super(pipelineRepository);
  }

  protected toEntity(model: PipelineWithPipelineToolsModel): PipelineEntity {
    return new PipelineEntity(model);
  }
}
