import { Injectable, Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { RunEntity } from "../runs/entities/run.entity";
import { RunsService } from "../runs/runs.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { RunToolDto } from "./dto/run-tool.dto";
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
  constructor(
    private toolsRepository: ToolRepository,
    private websocketsService: WebsocketsService,
    private runsService: RunsService
  ) {
    super(toolsRepository);
  }

  async run(
    orgname: string,
    toolId: string,
    runToolDto: RunToolDto
  ): Promise<RunEntity> {
    const tool = await this.toolsRepository.findOne(orgname, toolId);
    return this.runsService.runTool(orgname, tool, runToolDto);
  }

  protected toEntity(model: ToolModel): ToolEntity {
    return new ToolEntity(model);
  }
}
