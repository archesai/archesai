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
import { SearchQueryDto } from "../common/dto/search-query.dto";
import { CreateRunDto } from "./dto/create-run.dto";
import { RunEntity } from "./entities/run.entity";
import { RunsService } from "./runs.service";

@ApiBearerAuth()
@ApiTags("Runs")
@Controller("organizations/:orgname/runs")
export class RunsController
  implements BaseController<RunEntity, CreateRunDto, SearchQueryDto, any>
{
  constructor(private readonly runsService: RunsService) {}

  @ApiCrudOperation(Operation.CREATE, "run", RunEntity, true)
  @Post("/")
  async create(
    @Param("orgname") orgname: string,
    @Body() createRunDto: CreateRunDto
  ) {
    return this.runsService.create(orgname, createRunDto);
  }

  @ApiCrudOperation(Operation.FIND_ALL, "run", RunEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() searchQueryDto: SearchQueryDto
  ) {
    return this.runsService.findAll(orgname, searchQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "run", RunEntity, true)
  @Get("/:runId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("runId") runId: string
  ) {
    return this.runsService.findOne(orgname, runId);
  }

  @ApiCrudOperation(Operation.DELETE, "run", RunEntity, true)
  @Delete("/:runId")
  async remove(
    @Param("orgname") orgname: string,
    @Param("runId") runId: string
  ) {
    await this.runsService.remove(orgname, runId);
  }
}
