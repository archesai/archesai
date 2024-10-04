import { Injectable } from "@nestjs/common";
import { ApiToken } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { ApiTokenQueryDto } from "./dto/api-token-query.dto";
import { CreateApiTokenDto } from "./dto/create-api-token.dto";
import { UpdateApiTokenDto } from "./dto/update-api-token.dto";

@Injectable()
export class ApiTokenRepository
  implements
    BaseRepository<
      ApiToken,
      CreateApiTokenDto,
      ApiTokenQueryDto,
      UpdateApiTokenDto
    >
{
  constructor(private prisma: PrismaService) {}

  async create(
    orgname: string,
    createApiTokenDto: CreateApiTokenDto,
    snippet?: string
  ) {
    return this.prisma.apiToken.create({
      data: {
        chatbots: {
          connect: createApiTokenDto.chatbotIds.map((id) => ({ id })),
        },
        domains: createApiTokenDto.domains,
        key: snippet,
        name: createApiTokenDto.name,
        organization: {
          connect: {
            orgname,
          },
        },
        role: createApiTokenDto.role,
      },
      include: {
        chatbots: {
          select: {
            id: true,
            name: true,
          },
        },
      },
    });
  }

  async findAll(orgname: string, apiTokenQueryDto: ApiTokenQueryDto) {
    const count = await this.prisma.apiToken.count({
      where: { name: { contains: apiTokenQueryDto.name }, orgname },
    });
    const results = await this.prisma.apiToken.findMany({
      include: {
        chatbots: {
          select: {
            id: true,
            name: true,
          },
        },
      },
      orderBy: {
        [apiTokenQueryDto.sortBy]: apiTokenQueryDto.sortDirection,
      },
      skip: apiTokenQueryDto.offset,
      take: apiTokenQueryDto.limit,
      where: { name: { contains: apiTokenQueryDto.name }, orgname },
    });
    return { count, results };
  }

  async findOne(orgname: string, id: string) {
    return this.prisma.apiToken.findUniqueOrThrow({
      include: {
        chatbots: {
          select: {
            id: true,
            name: true,
          },
        },
      },
      where: { id },
    });
  }

  async remove(orgname: string, id: string) {
    await this.prisma.apiToken.delete({
      where: { id },
    });
  }

  async update(
    orgname: string,
    id: string,
    updateApiTokenDto: UpdateApiTokenDto
  ) {
    return this.prisma.apiToken.update({
      data: {
        ...(updateApiTokenDto.chatbotIds && {
          chatbots: {
            set: updateApiTokenDto.chatbotIds.map((id) => ({ id })),
          },
        }),
        ...(updateApiTokenDto.domains && {
          domains: updateApiTokenDto.domains,
        }),
        ...(updateApiTokenDto.name && { name: updateApiTokenDto.name }),
        ...(updateApiTokenDto.role && { role: updateApiTokenDto.role }),
      },
      include: {
        chatbots: {
          select: {
            id: true,
            name: true,
          },
        },
      },
      where: { id },
    });
  }
}
