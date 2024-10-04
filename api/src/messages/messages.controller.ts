import { Body, Controller, Get, Param, Post, Query, Req } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";
import { Request } from "express";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { PaginatedDto } from "../common/paginated.dto";
import { CreateMessageDto } from "./dto/create-message.dto";
import { MessageQueryDto } from "./dto/message-query.dto";
import { MessageEntity } from "./entities/message.entity";
import { MessagesService } from "./messages.service";

@ApiTags("Chatbots - Threads - Messages")
@ApiBearerAuth()
@Controller(
  "/organizations/:orgname/chatbots/:chatbotId/threads/:threadId/messages"
)
export class MessagesController {
  constructor(private readonly messagesService: MessagesService) {}

  @ApiCrudOperation(Operation.CREATE, "message", MessageEntity, true)
  @Post()
  async create(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string,
    @Param("threadId") threadId: string,
    @Body() createMessageDto: CreateMessageDto,
    @Req() req: Request
  ) {
    const controller = new AbortController();
    const handleRequestClose = () => {
      controller.abort();
      return;
    };
    // Listen for the aborted event on the request
    req.socket.on("close", handleRequestClose);
    const message = new MessageEntity(
      await this.messagesService.create(
        orgname,
        chatbotId,
        threadId,
        createMessageDto,
        controller.signal
      )
    );
    req.socket.off("close", handleRequestClose);
    return message;
  }

  @ApiCrudOperation(Operation.FIND_ALL, "message", MessageEntity, true)
  @Get()
  async findAll(
    @Query() messageQueryDto: MessageQueryDto,
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string,
    @Param("threadId") threadId: string
  ) {
    const { count, results } = await this.messagesService.findAll(
      orgname,
      chatbotId,
      threadId,
      messageQueryDto
    );
    return new PaginatedDto<MessageEntity>({
      metadata: {
        limit: messageQueryDto.limit,
        offset: messageQueryDto.offset,
        totalResults: count,
      },
      results: results.map((message) => new MessageEntity(message)),
    });
  }
}
