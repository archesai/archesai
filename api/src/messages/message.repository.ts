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
    answer: string
  ) {
    return this.prisma.message.create({
      data: {
        answer,
        question: createMessageDto.question,
        thread: {
          connect: {
            id: threadId,
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
