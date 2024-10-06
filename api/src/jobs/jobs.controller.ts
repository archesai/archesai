import { Controller, Delete, Get, Param, Query } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { PaginatedDto } from "../common/paginated.dto";
import { JobQueryDto } from "./dto/job-query.dto";
import { JobEntity } from "./entities/job.entity";
import { JobsService } from "./jobs.service";

@ApiBearerAuth()
@ApiTags("Organization - Jobs")
@Controller("/organizations/:orgname/jobs")
export class JobsController
  implements BaseController<JobEntity, undefined, JobQueryDto, undefined>
{
  constructor(private readonly jobsService: JobsService) {}

  @ApiCrudOperation(Operation.FIND_ALL, "job", JobEntity, true)
  @Get()
  async findAll(
    @Param("orgname") orgname: string,
    @Query() jobQueryDto: JobQueryDto
  ) {
    const { count, results } = await this.jobsService.findAll(
      orgname,
      jobQueryDto
    );
    return new PaginatedDto<JobEntity>({
      metadata: {
        limit: jobQueryDto.limit,
        offset: jobQueryDto.offset,
        totalResults: count,
      },
      results: results.map((chatbot) => new JobEntity(chatbot)),
    });
  }

  @ApiCrudOperation(Operation.GET, "job", JobEntity, true)
  @Get(":id")
  findOne(@Param("orgname") orgname: string, @Param("id") id: string) {
    return this.jobsService.findOne(orgname, id);
  }

  @ApiCrudOperation(Operation.DELETE, "job", JobEntity, true)
  @Delete(":id")
  remove(@Param("orgname") orgname: string, @Param("id") id: string) {
    return this.jobsService.remove(orgname, id);
  }
}
