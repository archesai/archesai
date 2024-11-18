import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Query,
} from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import { BaseController } from "../common/base.controller";
import {
  ApiCrudOperation,
  Operation,
} from "../common/decorators/api-crud-operation.decorator";
import { CreateRunDto } from "../common/dto/create-run.dto";
import { SearchQueryDto } from "../common/dto/search-query.dto";
import { ToolRunEntity } from "./entities/tool-run.entity";
import { ToolRunsService } from "./tool-runs.service";

@ApiBearerAuth()
@ApiTags("Tool Runs")
@Controller("organizations/:orgname/tool-runs")
export class ToolRunsController
  implements BaseController<ToolRunEntity, CreateRunDto, SearchQueryDto, any>
{
  constructor(private readonly toolRunsService: ToolRunsService) {}

  @ApiCrudOperation(Operation.CREATE, "tool run", ToolRunEntity, true)
  @Post("/")
  async create(
    @Param("orgname") orgname: string,
    @Body() createToolRunDto: CreateRunDto
  ) {
    return this.toolRunsService.create(orgname, createToolRunDto);
  }

  @ApiCrudOperation(Operation.FIND_ALL, "tool run", ToolRunEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() searchQueryDto: SearchQueryDto
  ) {
    return this.toolRunsService.findAll(orgname, searchQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "tool run", ToolRunEntity, true)
  @Get("/:toolRunId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("toolRunId") toolRunId: string
  ) {
    return this.toolRunsService.findOne(orgname, toolRunId);
  }

  @ApiCrudOperation(Operation.DELETE, "tool run", ToolRunEntity, true)
  @Delete("/:toolRunId")
  async remove(
    @Param("orgname") orgname: string,
    @Param("toolRunId") toolRunId: string
  ) {
    await this.toolRunsService.remove(orgname, toolRunId);
  }
}
