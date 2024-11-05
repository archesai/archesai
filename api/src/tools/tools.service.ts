import { Injectable, Logger } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { RunEntity } from "../runs/entities/run.entity";
import { RunsService } from "../runs/runs.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { RunToolDto } from "./dto/run-tool.dto";
import { ToolQueryDto } from "./dto/tool-query.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolEntity } from "./entities/tool.entity";
import { ToolRepository } from "./tool.repository";

@Injectable()
export class ToolsService
  implements
    BaseService<ToolEntity, CreateToolDto, ToolQueryDto, UpdateToolDto>
{
  private logger = new Logger(ToolsService.name);
  constructor(
    private toolsRepository: ToolRepository,
    private websocketsService: WebsocketsService,
    private runsService: RunsService
  ) {}

  async create(
    orgname: string,
    createToolDto: CreateToolDto
  ): Promise<ToolEntity> {
    const tool = await this.toolsRepository.create(orgname, createToolDto);
    return new ToolEntity(tool);
  }

  async findAll(orgname: string, toolsQueryDto: ToolQueryDto) {
    const { count, results } = await this.toolsRepository.findAll(
      orgname,
      toolsQueryDto
    );
    const toolEntities = results.map((tool) => new ToolEntity(tool));
    return new PaginatedDto<ToolEntity>({
      metadata: {
        limit: toolsQueryDto.limit,
        offset: toolsQueryDto.offset,
        totalResults: count,
      },
      results: toolEntities,
    });
  }

  async findOne(id: string) {
    return new ToolEntity(await this.toolsRepository.findOne(id));
  }

  async remove(orgname: string, toolId: string): Promise<void> {
    await this.toolsRepository.remove(orgname, toolId);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "tools"],
    });
  }

  async run(
    orgname: string,
    toolId: string,
    runToolDto: RunToolDto
  ): Promise<RunEntity> {
    const tool = await this.toolsRepository.findOne(toolId);
    return this.runsService.runTool(orgname, tool, runToolDto);
  }

  async update(orgname: string, id: string, updateToolDto: UpdateToolDto) {
    const tool = await this.toolsRepository.update(orgname, id, updateToolDto);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "tools"],
    });
    return new ToolEntity(tool);
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ToolUpdateInput) {
    const tool = await this.toolsRepository.updateRaw(orgname, id, raw);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "tools"],
    });
    return new ToolEntity(tool);
  }
}
