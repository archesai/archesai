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
} from "../common/api-crud-operation.decorator";
import { PaginatedDto } from "../common/paginated.dto";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadAggregates } from "./dto/thread-aggregates.dto";
import { ThreadQueryDto } from "./dto/thread-query.dto";
import { ThreadEntity } from "./entities/thread.entity";
import { ThreadsService } from "./threads.service";

@ApiBearerAuth()
@ApiTags("Threads")
@Controller("/organizations/:orgname/threads")
export class ThreadsController {
  constructor(private readonly threadsService: ThreadsService) {}

  @ApiCrudOperation(Operation.CREATE, "thread", ThreadEntity, false)
  @Post()
  async create(
    @Param("orgname") orgname: string,
    @Body() createThreadDto: CreateThreadDto
  ) {
    return new ThreadEntity(
      await this.threadsService.create(orgname, createThreadDto)
    );
  }

  @ApiCrudOperation(
    Operation.FIND_ALL,
    "thread",
    ThreadEntity,
    true,
    ThreadAggregates
  )
  @Get()
  async findAll(
    @Query() threadQueryDto: ThreadQueryDto,
    @Param("orgname") orgname: string
  ) {
    const { aggregates, count, threads } = await this.threadsService.findAll(
      orgname,
      threadQueryDto
    );
    return new PaginatedDto<ThreadEntity, ThreadAggregates>({
      aggregates,
      metadata: {
        limit: threadQueryDto.limit,
        offset: threadQueryDto.offset,
        totalResults: count,
      },
      results: threads.map((thread) => new ThreadEntity(thread)),
    });
  }

  @ApiCrudOperation(Operation.GET, "thread", ThreadEntity, true)
  @Get(":threadId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("threadId") threadId: string
  ) {
    return new ThreadEntity(
      await this.threadsService.findOne(orgname, threadId)
    );
  }

  @ApiCrudOperation(Operation.DELETE, "thread", ThreadEntity, true)
  @Delete(":threadId")
  async remove(
    @Param("orgname") orgname: string,
    @Param("threadId") threadId: string
  ) {
    return this.threadsService.remove(orgname, threadId);
  }
}
