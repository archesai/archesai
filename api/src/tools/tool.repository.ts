import { Injectable } from "@nestjs/common";
import { Prisma, Tool } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { ToolQueryDto } from "./dto/tool-query.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";

@Injectable()
export class ToolRepository
  implements BaseRepository<Tool, CreateToolDto, ToolQueryDto, UpdateToolDto>
{
  constructor(private prisma: PrismaService) {}

  async create(orgname: string, createToolDto: CreateToolDto) {
    return this.prisma.tool.create({
      data: {
        organization: {
          connect: {
            orgname,
          },
        },
        ...createToolDto,
      },
    });
  }

  async findAll(orgname: string, toolQueryDto: ToolQueryDto) {
    const whereConditions = {
      createdAt: {
        gte: toolQueryDto.startDate,
        lte: toolQueryDto.endDate,
      },
      orgname,
    };
    if (toolQueryDto.filters) {
      toolQueryDto.filters.forEach((filter) => {
        whereConditions[filter.field] = { contains: filter.value };
      });
    }
    const count = await this.prisma.tool.count({
      where: whereConditions,
    });
    const tool = await this.prisma.tool.findMany({
      orderBy: {
        [toolQueryDto.sortBy]: toolQueryDto.sortDirection,
      },
      skip: toolQueryDto.offset,
      take: toolQueryDto.limit,
      where: whereConditions,
    });
    return { count, results: tool };
  }

  async findOne(id: string) {
    return this.prisma.tool.findUniqueOrThrow({
      where: { id },
    });
  }

  async remove(orgname: string, id: string) {
    await this.prisma.tool.delete({
      where: { id },
    });
  }

  async update(orgname: string, id: string, updateToolDto: UpdateToolDto) {
    return this.prisma.tool.update({
      data: {
        ...updateToolDto,
      },
      where: {
        id,
      },
    });
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ToolUpdateInput) {
    return this.prisma.tool.update({
      data: raw,
      where: {
        id,
      },
    });
  }
}
