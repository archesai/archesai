import { Injectable, Logger } from "@nestjs/common";
import { Prisma, Tool } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { RunEntity } from "../runs/entities/run.entity";
import { RunsService } from "../runs/runs.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { RunToolDto } from "./dto/run-tool.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolEntity } from "./entities/tool.entity";
import { ToolRepository } from "./tool.repository";

@Injectable()
export class ToolsService extends BaseService<
  ToolEntity,
  CreateToolDto,
  UpdateToolDto,
  ToolRepository,
  Tool,
  Prisma.ToolInclude,
  Prisma.ToolSelect
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

  protected toEntity(model: Tool): ToolEntity {
    return new ToolEntity(model);
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ToolUpdateInput) {
    const tool = await this.toolsRepository.updateRaw(orgname, id, raw);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "tools"],
    });
    return this.toEntity(tool);
  }
}
