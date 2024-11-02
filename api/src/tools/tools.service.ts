import { Inject, Injectable, Logger } from "@nestjs/common";
import { Prisma, Tool } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { ToolQueryDto } from "./dto/tool-query.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolEntity } from "./entities/tool.entity";
import { ToolRepository } from "./tool.repository";

@Injectable()
export class ToolsService
  implements BaseService<Tool, CreateToolDto, ToolQueryDto, UpdateToolDto>
{
  private logger = new Logger(ToolsService.name);
  constructor(
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private toolsRepository: ToolRepository,
    private websocketsService: WebsocketsService
  ) {}

  async create(
    orgname: string,
    createToolDto: CreateToolDto
  ): Promise<ToolEntity> {
    const tool = await this.toolsRepository.create(orgname, createToolDto);
    return new ToolEntity(tool);
  }

  async findAll(orgname: string, toolsQueryDto: ToolQueryDto) {
    return this.toolsRepository.findAll(orgname, toolsQueryDto);
  }

  async findOne(id: string) {
    const tool = await this.toolsRepository.findOne(id);
    const populated = await this.populateReadUrl(tool);
    return populated;
  }

  async populateReadUrl(tool: Tool): Promise<Tool> {
    if (
      tool.url?.startsWith(
        `https://storage.googleapis.com/archesai/storage/${tool.orgname}/`
      )
    ) {
      const path = tool.url
        .replace(
          `https://storage.googleapis.com/archesai/storage/${tool.orgname}/`,
          ""
        )
        .split("?")[0];

      try {
        const read = await this.storageService.getSignedUrl(
          tool.orgname,
          decodeURIComponent(path),
          "read"
        );
        tool.url = read;
      } catch (e) {
        this.logger.warn(e);
        tool.url = "";
      }
    }

    return tool;
  }

  async remove(orgname: string, toolsId: string): Promise<void> {
    await this.toolsRepository.remove(orgname, toolsId);
    this.websocketsService.socket.to(orgname).emit("update");
  }

  async update(orgname: string, id: string, updateToolDto: UpdateToolDto) {
    const tools = await this.toolsRepository.update(orgname, id, updateToolDto);
    this.websocketsService.socket.to(orgname).emit("update");
    return tools;
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ToolUpdateInput) {
    const tools = await this.toolsRepository.updateRaw(orgname, id, raw);
    this.websocketsService.socket.to(orgname).emit("update");
    return tools;
  }
}
