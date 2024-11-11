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
} from "../common/decorators/api-crud-operation.decorator";
import { SearchQueryDto } from "../common/dto/search-query.dto";
import { CreateLabelDto } from "./dto/create-label.dto";
import { UpdateLabelDto } from "./dto/update-label.dto";
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

  @ApiCrudOperation(Operation.FIND_ALL, "label", LabelEntity, true)
  @Get()
  async findAll(
    @Query() searchQueryDto: SearchQueryDto,
    @Param("orgname") orgname: string
  ) {
    return this.labelsService.findAll(orgname, searchQueryDto);
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

  @ApiCrudOperation(Operation.UPDATE, "label", LabelEntity, true)
  @Patch(":labelId")
  async update(
    @Param("orgname") orgname: string,
    @Param("labelId") labelId: string,
    @Body() updateLabelDto: UpdateLabelDto
  ) {
    return this.labelsService.update(orgname, labelId, updateLabelDto);
  }
}
