import { Injectable } from "@nestjs/common";
import { Chatbot } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { OrganizationsService } from "../organizations/organizations.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { ChatbotRepository } from "./chatbot.repository";
import { ChatbotQueryDto } from "./dto/chatbot-query.dto";
import { CreateChatbotDto } from "./dto/create-chatbot.dto";
import { UpdateChatbotDto } from "./dto/update-chatbot.dto";

@Injectable()
export class ChatbotsService
  implements
    BaseService<Chatbot, CreateChatbotDto, ChatbotQueryDto, UpdateChatbotDto>
{
  constructor(
    private chatbotRepository: ChatbotRepository,
    private websocketsService: WebsocketsService,
    private organizationsService: OrganizationsService
  ) {}

  async create(orgname: string, createChatbotDto: CreateChatbotDto) {
    const chatbot = await this.chatbotRepository.create(
      orgname,
      createChatbotDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return chatbot;
  }

  async findAll(orgname: string, chatbotQueryDto: ChatbotQueryDto) {
    return this.chatbotRepository.findAll(orgname, chatbotQueryDto);
  }

  async findOne(chatbotId: string) {
    return this.chatbotRepository.findOne(chatbotId);
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
    const chatbot = await this.chatbotRepository.update(
      orgname,
      chatbotId,
      updateChatbotDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return chatbot;
  }
}
