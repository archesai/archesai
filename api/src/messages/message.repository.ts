import { Injectable } from "@nestjs/common";
import { Message, Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateMessageDto } from "./dto/create-message.dto";

@Injectable()
export class MessageRepository extends BaseRepository<
  Message,
  CreateMessageDto,
  undefined,
  Prisma.MessageInclude,
  Prisma.MessageSelect,
  Prisma.MessageUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.message);
  }

  async create(
    orgname: string,
    createMessageDto: CreateMessageDto,
    additionalData: {
      answer: string;
      threadId: string;
    }
  ) {
    return this.prisma.message.create({
      data: {
        answer: additionalData.answer,
        organization: {
          connect: {
            orgname,
          },
        },
        question: createMessageDto.question,
        thread: {
          connect: {
            id: additionalData.threadId,
          },
        },
      },
    });
  }
}
