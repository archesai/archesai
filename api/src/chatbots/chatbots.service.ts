import { BadRequestException, Injectable } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { WebsocketsService } from "../websockets/websockets.service";
import { ChatbotRepository } from "./chatbot.repository";
import { ChatbotQueryDto } from "./dto/chatbot-query.dto";
import { CreateChatbotDto } from "./dto/create-chatbot.dto";
import { UpdateChatbotDto } from "./dto/update-chatbot.dto";
import { ChatbotEntity } from "./entities/chatbot.entity";

@Injectable()
export class ChatbotsService
  implements
    BaseService<
      ChatbotEntity,
      CreateChatbotDto,
      ChatbotQueryDto,
      UpdateChatbotDto
    >
{
  constructor(
    private chatbotRepository: ChatbotRepository,
    private websocketsService: WebsocketsService
  ) {}

  async create(orgname: string, createChatbotDto: CreateChatbotDto) {
    const chatbot = await this.chatbotRepository.create(
      orgname,
      createChatbotDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return new ChatbotEntity(chatbot);
  }

  async findAll(orgname: string, chatbotQueryDto: ChatbotQueryDto) {
    const { count, results } = await this.chatbotRepository.findAll(
      orgname,
      chatbotQueryDto
    );
    const chatbotEntities = results.map(
      (chatbot) => new ChatbotEntity(chatbot)
    );
    return new PaginatedDto<ChatbotEntity>({
      metadata: {
        limit: chatbotQueryDto.limit,
        offset: chatbotQueryDto.offset,
        totalResults: count,
      },
      results: chatbotEntities,
    });
  }

  async findOne(chatbotId: string) {
    return new ChatbotEntity(await this.chatbotRepository.findOne(chatbotId));
  }

  async remove(orgname: string, chatbotId: string) {
    await this.chatbotRepository.remove(orgname, chatbotId);
    this.websocketsService.socket.to(orgname).emit("update");
  }

  async update(
    orgname: string,
    chatbotId: string,
    updateChatbotDto: UpdateChatbotDto
  ) {
    if (
      updateChatbotDto.llmBase &&
      !["gpt-4o", "gpt-4o-mini"].includes(updateChatbotDto.llmBase)
    ) {
      throw new BadRequestException("Invalid LLM base");
    }
    const chatbot = await this.chatbotRepository.update(
      orgname,
      chatbotId,
      updateChatbotDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return new ChatbotEntity(chatbot);
  }
}
