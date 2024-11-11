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

import { BaseController } from "../common/base.controller";
import {
  ApiCrudOperation,
  Operation,
} from "../common/decorators/api-crud-operation.decorator";
import { SearchQueryDto } from "../common/dto/search-query.dto";
import { CreateToolDto } from "./dto/create-tool.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolEntity } from "./entities/tool.entity";
import { ToolsService } from "./tools.service";

@ApiBearerAuth()
@ApiTags("Tools")
@Controller("organizations/:orgname/tools")
export class ToolsController
  implements
    BaseController<ToolEntity, CreateToolDto, SearchQueryDto, UpdateToolDto>
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
    @Query() searchQueryDto: SearchQueryDto
  ) {
    return this.toolsService.findAll(orgname, searchQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "tool", ToolEntity, true)
  @Get("/:toolId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string
  ) {
    return this.toolsService.findOne(orgname, toolId);
  }

  @ApiCrudOperation(Operation.DELETE, "tool", ToolEntity, true)
  @Delete("/:toolId")
  async remove(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string
  ) {
    await this.toolsService.remove(orgname, toolId);
  }

  @ApiCrudOperation(Operation.UPDATE, "tool", ToolEntity, true)
  @Patch("/:toolId")
  async update(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string,
    @Body() updateContentDto: UpdateToolDto
  ) {
    return this.toolsService.update(orgname, toolId, updateContentDto);
  }
}
