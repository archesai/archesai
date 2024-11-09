import { Injectable, Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolEntity, ToolModel } from "./entities/tool.entity";
import { ToolRepository } from "./tool.repository";

@Injectable()
export class ToolsService extends BaseService<
  ToolEntity,
  CreateToolDto,
  UpdateToolDto,
  ToolRepository,
  ToolModel
> {
  private logger = new Logger(ToolsService.name);
  constructor(private toolsRepository: ToolRepository) {
    super(toolsRepository);
  }

  protected toEntity(model: ToolModel): ToolEntity {
    return new ToolEntity(model);
  }
}
