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

import {
  ApiCrudOperation,
  Operation,
} from "../common/decorators/api-crud-operation.decorator";
import { CreateLabelDto } from "./dto/create-label.dto";
import { LabelAggregates } from "./dto/label-aggregates.dto";
import { LabelQueryDto } from "./dto/label-query.dto";
import { LabelEntity } from "./entities/label.entity";
import { LabelsService } from "./labels.service";

@ApiBearerAuth()
@ApiTags("Labels")
@Controller("/organizations/:orgname/labels")
export class LabelsController {
  constructor(private readonly labelsService: LabelsService) {}

  @ApiCrudOperation(Operation.CREATE, "label", LabelEntity, false)
  @Post()
  async create(
    @Param("orgname") orgname: string,
    @Body() createLabelDto: CreateLabelDto
  ) {
    return this.labelsService.create(orgname, createLabelDto);
  }

  @ApiCrudOperation(
    Operation.FIND_ALL,
    "label",
    LabelEntity,
    true,
    LabelAggregates
  )
  @Get()
  async findAll(
    @Query() labelQueryDto: LabelQueryDto,
    @Param("orgname") orgname: string
  ) {
    return this.labelsService.findAll(orgname, labelQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "label", LabelEntity, true)
  @Get(":labelId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("labelId") labelId: string
  ) {
    return this.labelsService.findOne(orgname, labelId);
  }

  @ApiCrudOperation(Operation.DELETE, "label", LabelEntity, true)
  @Delete(":labelId")
  async remove(
    @Param("orgname") orgname: string,
    @Param("labelId") labelId: string
  ) {
    return this.labelsService.remove(orgname, labelId);
  }
}
