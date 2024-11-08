import { Controller, Get, Param, Query } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { SearchQueryDto } from "../common/search-query";
import { RunEntity } from "./entities/run.entity";
import { RunDetailedEntity } from "./entities/run-detailed.entity";
import { RunsService } from "./runs.service";

@ApiBearerAuth()
@ApiTags("Runs")
@Controller("/organizations/:orgname/runs")
export class RunsController
  implements BaseController<RunEntity, undefined, SearchQueryDto, undefined>
{
  constructor(private readonly runsService: RunsService) {}

  @ApiCrudOperation(Operation.FIND_ALL, "run", RunEntity, true)
  @Get()
  async findAll(
    @Param("orgname") orgname: string,
    @Query() searchQuery: SearchQueryDto
  ) {
    return this.runsService.findAll(orgname, searchQuery);
  }

  @ApiCrudOperation(Operation.GET, "run", RunDetailedEntity, true)
  @Get(":id")
  findOne(@Param("orgname") orgname: string, @Param("id") id: string) {
    return this.runsService.findOne(orgname, id);
  }
}
