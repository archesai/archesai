import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  Query,
} from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { RunEntity } from "../runs/entities/run.entity";
import { CreateToolDto } from "./dto/create-tool.dto";
import { RunToolDto } from "./dto/run-tool.dto";
import { ToolQueryDto } from "./dto/tool-query.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolEntity } from "./entities/tool.entity";
import { ToolsService } from "./tools.service";

@ApiBearerAuth()
@ApiTags("Tools")
@Controller("organizations/:orgname/tools")
export class ToolsController
  implements
    BaseController<ToolEntity, CreateToolDto, ToolQueryDto, UpdateToolDto>
{
  constructor(private readonly toolsService: ToolsService) {}

  @ApiCrudOperation(Operation.CREATE, "tool", ToolEntity, true)
  @Post("/")
  async create(
    @Param("orgname") orgname: string,
    @Body() createToolDto: CreateToolDto
  ) {
    return this.toolsService.create(orgname, createToolDto);
  }

  @ApiCrudOperation(Operation.FIND_ALL, "tool", ToolEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() toolQueryDto: ToolQueryDto
  ) {
    return this.toolsService.findAll(orgname, toolQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "tool", ToolEntity, true)
  @Get("/:toolId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string
  ) {
    return this.toolsService.findOne(toolId);
  }

  @ApiCrudOperation(Operation.DELETE, "tools", ToolEntity, true)
  @Delete("/:toolId")
  remove(@Param("orgname") orgname: string, @Param("toolId") toolId: string) {
    return this.toolsService.remove(orgname, toolId);
  }

  @ApiCrudOperation(Operation.CREATE, "run", RunEntity, true)
  @Post("/:toolId/run")
  async run(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string,
    @Body() runToolDto: RunToolDto
  ) {
    return this.toolsService.run(orgname, toolId, runToolDto);
  }

  @ApiCrudOperation(Operation.UPDATE, "tools", ToolEntity, true)
  @Patch("/:toolId")
  async update(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string,
    @Body() updateContentDto: UpdateToolDto
  ) {
    return this.toolsService.update(orgname, toolId, updateContentDto);
  }
}
