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

import { CrudOperation, Operation } from "../common/crud-operation.decorator";
import { PaginatedDto } from "../common/paginated.dto";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadAggregates } from "./dto/thread-aggregates.dto";
import { ThreadQueryDto } from "./dto/thread-query.dto";
import { ThreadEntity } from "./entities/thread.entity";
import { ThreadsService } from "./threads.service";

@ApiBearerAuth()
@ApiTags("Chatbots - Threads")
@Controller("/organizations/:orgname/chatbots/:chatbotId/threads")
export class ThreadsController {
  constructor(private readonly threadsService: ThreadsService) {}

  @CrudOperation(Operation.CREATE, "thread", ThreadEntity, false)
  @Post()
  async create(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string,
    @Body() createThreadDto: CreateThreadDto
  ) {
    return new ThreadEntity(
      await this.threadsService.create(orgname, chatbotId, createThreadDto)
    );
  }

  @CrudOperation(
    Operation.FIND_ALL,
    "thread",
    ThreadEntity,
    true,
    ThreadAggregates
  )
  @Get()
  async findAll(
    @Query() threadQueryDto: ThreadQueryDto,
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string
  ) {
    const { aggregates, count, threads } = await this.threadsService.findAll(
      orgname,
      chatbotId,
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

  @CrudOperation(Operation.GET, "thread", ThreadEntity, true)
  @Get(":threadId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string,
    @Param("threadId") threadId: string
  ) {
    return new ThreadEntity(
      await this.threadsService.findOne(orgname, chatbotId, threadId)
    );
  }

  @CrudOperation(Operation.DELETE, "thread", ThreadEntity, true)
  @Delete(":threadId")
  async remove(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string,
    @Param("threadId") threadId: string
  ) {
    return this.threadsService.remove(orgname, chatbotId, threadId);
  }
}
