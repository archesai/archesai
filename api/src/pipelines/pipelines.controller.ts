import {
  Body,
  Controller,
  Delete,
  Get,
  Logger,
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
import { SearchQueryDto } from "../common/search-query";
import { CreatePipelineDto } from "./dto/create-pipeline.dto";
import { UpdatePipelineDto } from "./dto/update-pipeline.dto";
import { PipelineEntity } from "./entities/pipeline.entity";
import { PipelinesService } from "./pipelines.service";

@ApiBearerAuth()
@ApiTags("Pipelines")
@Controller("organizations/:orgname/pipelines")
export class PipelinesController
  implements
    BaseController<
      PipelineEntity,
      CreatePipelineDto,
      SearchQueryDto,
      UpdatePipelineDto
    >
{
  private logger = new Logger(PipelinesController.name);

  constructor(private readonly pipelinesService: PipelinesService) {}

  @ApiCrudOperation(Operation.CREATE, "pipeline", PipelineEntity, true)
  @Post("/")
  async create(
    @Param("orgname") orgname: string,
    @Body() createPipelineDto: CreatePipelineDto
  ) {
    this.logger.log(JSON.stringify(createPipelineDto, null, 2));
    return this.pipelinesService.create(orgname, createPipelineDto);
  }

  @ApiCrudOperation(Operation.FIND_ALL, "pipeline", PipelineEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() searchQueryDto: SearchQueryDto
  ) {
    return this.pipelinesService.findAll(orgname, searchQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "pipeline", PipelineEntity, true)
  @Get("/:pipelineId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("pipelineId") pipelineId: string
  ) {
    return this.pipelinesService.findOne(orgname, pipelineId);
  }

  @ApiCrudOperation(Operation.DELETE, "pipeline", PipelineEntity, true)
  @Delete("/:pipelineId")
  remove(
    @Param("orgname") orgname: string,
    @Param("pipelineId") pipelineId: string
  ) {
    return this.pipelinesService.remove(orgname, pipelineId);
  }

  @ApiCrudOperation(Operation.UPDATE, "pipeline", PipelineEntity, true)
  @Patch("/:pipelineId")
  async update(
    @Param("orgname") orgname: string,
    @Param("pipelineId") pipelineId: string,
    @Body() updatePipelineDto: UpdatePipelineDto
  ) {
    return this.pipelinesService.update(orgname, pipelineId, updatePipelineDto);
  }
}
