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
    snippet?: string,
    id?: string,
    uid?: string
  ) {
    return this.prisma.apiToken.create({
      data: {
        domains: createApiTokenDto.domains,
        id,
        key: snippet,
        name: createApiTokenDto.name,
        organization: {
          connect: {
            orgname,
          },
        },
        role: createApiTokenDto.role,
        user: {
          connect: {
            id: uid,
          },
        },
      },
    });
  }

  async findAll(orgname: string, apiTokenQueryDto: ApiTokenQueryDto) {
    const whereConditions = {
      createdAt: {
        gte: apiTokenQueryDto.startDate,
        lte: apiTokenQueryDto.endDate,
      },
      orgname,
    };
    if (apiTokenQueryDto.filters) {
      apiTokenQueryDto.filters.forEach((filter) => {
        whereConditions[filter.field] = { contains: filter.value };
      });
    }

    const count = await this.prisma.apiToken.count({
      where: whereConditions,
    });
    const results = await this.prisma.apiToken.findMany({
      orderBy: {
        [apiTokenQueryDto.sortBy]: apiTokenQueryDto.sortDirection,
      },
      skip: apiTokenQueryDto.offset,
      take: apiTokenQueryDto.limit,
      where: whereConditions,
    });
    return { count, results };
  }

  async findOne(orgname: string, id: string) {
    return this.prisma.apiToken.findUniqueOrThrow({
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
        ...updateApiTokenDto,
      },
      where: { id },
    });
  }
}
