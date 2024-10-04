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

import { Roles } from "../auth/decorators/roles.decorator";
import { BaseController } from "../common/base.controller";
import { CrudOperation, Operation } from "../common/crud-operation.decorator";
import { PaginatedDto } from "../common/paginated.dto";
import { ChatbotsService } from "./chatbots.service";
import { ChatbotQueryDto } from "./dto/chatbot-query.dto";
import { CreateChatbotDto } from "./dto/create-chatbot.dto";
import { UpdateChatbotDto } from "./dto/update-chatbot.dto";
import { ChatbotEntity } from "./entities/chatbot.entity";

@Roles("ADMIN")
@ApiBearerAuth()
@ApiTags("Chatbots")
@Controller("/organizations/:orgname/chatbots")
export class ChatbotsController
  implements
    BaseController<
      ChatbotEntity,
      CreateChatbotDto,
      ChatbotQueryDto,
      UpdateChatbotDto
    >
{
  constructor(private readonly chatbotsService: ChatbotsService) {}

  @CrudOperation(Operation.CREATE, "chatbot", ChatbotEntity, true)
  @Post("/")
  async create(
    @Param("orgname") orgname: string,
    @Body() createChatbotDto: CreateChatbotDto
  ) {
    return new ChatbotEntity(
      await this.chatbotsService.create(orgname, createChatbotDto)
    );
  }

  @CrudOperation(Operation.FIND_ALL, "chatbot", ChatbotEntity, true)
  @Get()
  async findAll(
    @Param("orgname") orgname: string,
    @Query() chatbotQueryDto: ChatbotQueryDto
  ) {
    const { count, results } = await this.chatbotsService.findAll(
      orgname,
      chatbotQueryDto
    );
    return new PaginatedDto<ChatbotEntity>({
      metadata: {
        limit: chatbotQueryDto.limit,
        offset: chatbotQueryDto.offset,
        totalResults: count,
      },
      results: results.map((chatbot) => new ChatbotEntity(chatbot)),
    });
  }

  @CrudOperation(Operation.GET, "chatbot", ChatbotEntity, true)
  @Get(":chatbotId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string
  ) {
    return new ChatbotEntity(await this.chatbotsService.findOne(chatbotId));
  }

  @Delete(":chatbotId")
  @CrudOperation(Operation.DELETE, "chatbot", ChatbotEntity, true)
  async remove(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string
  ) {
    return this.chatbotsService.remove(orgname, chatbotId);
  }

  @CrudOperation(Operation.UPDATE, "chatbot", ChatbotEntity, true)
  @Patch(":chatbotId")
  async update(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string,
    @Body() updateChatbotDto: UpdateChatbotDto
  ) {
    return new ChatbotEntity(
      await this.chatbotsService.update(orgname, chatbotId, updateChatbotDto)
    );
  }
}
