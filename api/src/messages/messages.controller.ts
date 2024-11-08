import { Body, Controller, Get, Param, Post, Query, Req } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";
import { Request } from "express";

import { BaseController } from "../common/base.controller";
import {
  ApiCrudOperation,
  Operation,
} from "../common/decorators/api-crud-operation.decorator";
import { SearchQueryDto } from "../common/dto/search-query.dto";
import { CreateMessageDto } from "./dto/create-message.dto";
import { MessageEntity } from "./entities/message.entity";
import { MessagesService } from "./messages.service";

@ApiTags("Chatbots - Threads - Messages")
@ApiBearerAuth()
@Controller("/organizations/:orgname/threads/:threadId/messages")
export class MessagesController
  implements
    BaseController<MessageEntity, CreateMessageDto, SearchQueryDto, undefined>
{
  constructor(private readonly messagesService: MessagesService) {}

  @ApiCrudOperation(Operation.CREATE, "message", MessageEntity, true)
  @Post()
  async create(
    @Param("orgname") orgname: string,
    @Body() createMessageDto: CreateMessageDto,
    @Param("threadId") threadId: string,
    @Req() req: Request
  ) {
    const controller = new AbortController();
    const handleRequestClose = () => {
      controller.abort();
      return;
    };
    // Listen for the aborted event on the request
    req.socket.on("close", handleRequestClose);
    const message = await this.messagesService.create(
      orgname,
      createMessageDto,
      { threadId }
    );

    req.socket.off("close", handleRequestClose);
    return message;
  }

  @ApiCrudOperation(Operation.FIND_ALL, "message", MessageEntity, true)
  @Get()
  async findAll(
    @Param("orgname") orgname: string,
    @Query() searchQueryDto: SearchQueryDto,
    @Param("threadId") threadId: string
  ) {
    return this.messagesService.findAll(orgname, {
      filters: [
        {
          field: "threadId",
          operator: "equals",
          value: threadId,
        },
        ...searchQueryDto.filters,
      ],
      ...searchQueryDto,
    });
  }
}
