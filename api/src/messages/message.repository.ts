import { Injectable } from "@nestjs/common";

import { MessageQueryDto } from "../messages/dto/message-query.dto";
import { PrismaService } from "../prisma/prisma.service";
import { CreateMessageDto } from "./dto/create-message.dto";

@Injectable()
export class MessageRepository {
  constructor(private prisma: PrismaService) {}

  async create(
    threadId: string,
    createMessageDto: CreateMessageDto,
    answer: string,
    credits: number,
    citations: {
      contentId: string;
      similarity: number;
      text: string;
    }[]
  ) {
    return this.prisma.message.create({
      data: {
        answer,
        answerLength: createMessageDto.answerLength,
        citations: {
          createMany: {
            data: citations.map((citation) => ({
              contentId: citation.contentId,
              similarity: citation.similarity,
            })),
          },
        },
        contextLength: createMessageDto.contextLength,
        credits,
        question: createMessageDto.question,
        temperature: createMessageDto.temperature,
        thread: {
          connect: {
            id: threadId,
          },
        },
        topK: createMessageDto.topK,
      },
      include: {
        citations: {
          include: {
            message: true,
          },
        },
      },
    });
  }

  async findAll(
    orgname: string,
    threadId: string,
    messageQueryDto: MessageQueryDto
  ) {
    const count = await this.prisma.message.count({
      where: { threadId },
    });
    const messages = await this.prisma.message.findMany({
      include: {
        citations: {
          include: {
            message: true,
          },
        },
      },
      orderBy: {
        [messageQueryDto.sortBy]: messageQueryDto.sortDirection,
      },
      skip: messageQueryDto.offset,
      take: messageQueryDto.limit,
      where: { threadId },
    });
    return { count, results: messages };
  }
}
