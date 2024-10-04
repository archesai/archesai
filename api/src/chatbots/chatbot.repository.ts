import { Injectable } from "@nestjs/common";
import { Chatbot } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { ChatbotQueryDto } from "./dto/chatbot-query.dto";
import { CreateChatbotDto } from "./dto/create-chatbot.dto";
import { UpdateChatbotDto } from "./dto/update-chatbot.dto";

@Injectable()
export class ChatbotRepository
  implements
    BaseRepository<
      Chatbot,
      CreateChatbotDto,
      ChatbotQueryDto,
      UpdateChatbotDto
    >
{
  constructor(private prisma: PrismaService) {}

  async create(orgname: string, createChatbotDto: CreateChatbotDto) {
    return this.prisma.chatbot.create({
      data: {
        description: createChatbotDto.description,

        name: createChatbotDto.name,
        organization: {
          connect: {
            orgname,
          },
        },
        ...(createChatbotDto.llmBase && {
          llmBase: createChatbotDto.llmBase,
        }),
      },
    });
  }

  async findAll(orgname: string, chatbotQueryDto: ChatbotQueryDto) {
    const count = await this.prisma.chatbot.count({
      where: { name: { contains: chatbotQueryDto.name }, orgname },
    });
    const results = await this.prisma.chatbot.findMany({
      orderBy: {
        [chatbotQueryDto.sortBy]: chatbotQueryDto.sortDirection,
      },
      skip: chatbotQueryDto.offset,
      take: chatbotQueryDto.limit,
      where: { name: { contains: chatbotQueryDto.name }, orgname },
    });
    return { count, results };
  }

  async findOne(id: string) {
    return this.prisma.chatbot.findUniqueOrThrow({
      where: { id },
    });
  }

  async remove(orgname: string, id: string) {
    await this.prisma.chatbot.delete({
      where: { id },
    });
  }

  async update(
    orgname: string,
    chatbotId: string,
    updateChatbotDto: UpdateChatbotDto
  ) {
    return this.prisma.chatbot.update({
      data: {
        ...updateChatbotDto,
      },
      where: { id: chatbotId },
    });
  }

  async updateChatbotName(chatbotId: string, name: string) {
    return this.prisma.chatbot.update({
      data: {
        name,
      },
      where: { id: chatbotId },
    });
  }
}
