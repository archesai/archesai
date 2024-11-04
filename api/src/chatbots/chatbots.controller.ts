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
import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { ChatbotsService } from "./chatbots.service";
import { ChatbotQueryDto } from "./dto/chatbot-query.dto";
import { CreateChatbotDto } from "./dto/create-chatbot.dto";
import { UpdateChatbotDto } from "./dto/update-chatbot.dto";
import { ChatbotEntity } from "./entities/chatbot.entity";

@ApiBearerAuth()
@Roles("ADMIN")
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

  @ApiCrudOperation(Operation.CREATE, "chatbot", ChatbotEntity, true)
  @Post("/")
  async create(
    @Param("orgname") orgname: string,
    @Body() createChatbotDto: CreateChatbotDto
  ) {
    return new ChatbotEntity(
      await this.chatbotsService.create(orgname, createChatbotDto)
    );
  }

  @ApiCrudOperation(Operation.FIND_ALL, "chatbot", ChatbotEntity, true)
  @Get()
  async findAll(
    @Param("orgname") orgname: string,
    @Query() chatbotQueryDto: ChatbotQueryDto
  ) {
    return this.chatbotsService.findAll(orgname, chatbotQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "chatbot", ChatbotEntity, true)
  @Get(":chatbotId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string
  ) {
    return this.chatbotsService.findOne(chatbotId);
  }

  @Delete(":chatbotId")
  @ApiCrudOperation(Operation.DELETE, "chatbot", ChatbotEntity, true)
  async remove(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string
  ) {
    return this.chatbotsService.remove(orgname, chatbotId);
  }

  @ApiCrudOperation(Operation.UPDATE, "chatbot", ChatbotEntity, true)
  @Patch(":chatbotId")
  async update(
    @Param("orgname") orgname: string,
    @Param("chatbotId") chatbotId: string,
    @Body() updateChatbotDto: UpdateChatbotDto
  ) {
    return this.chatbotsService.update(orgname, chatbotId, updateChatbotDto);
  }
}
