import { Injectable } from "@nestjs/common";
import { Content, Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { ContentQueryDto } from "./dto/content-query.dto";
import { CreateContentDto } from "./dto/create-content.dto";
import { UpdateContentDto } from "./dto/update-content.dto";

@Injectable()
export class ContentRepository
  implements
    BaseRepository<
      Content,
      CreateContentDto,
      ContentQueryDto,
      UpdateContentDto
    >
{
  constructor(private prisma: PrismaService) {}

  async create(orgname: string, createContentDto: CreateContentDto) {
    return this.prisma.content.create({
      data: {
        name: createContentDto.name,
        organization: {
          connect: {
            orgname,
          },
        },
        url: createContentDto.url,
      },
      include: {
        textChunks: true,
      },
    });
  }

  async findAll(orgname: string, contentQueryDto: ContentQueryDto) {
    const count = await this.prisma.content.count({
      where: {
        createdAt: {
          gte: contentQueryDto.startDate,
          lte: contentQueryDto.endDate,
        },
        orgname,
        ...(contentQueryDto.searchTerm
          ? {
              OR: [{ name: { contains: contentQueryDto.searchTerm } }],
            }
          : undefined),
      },
    });
    const content = await this.prisma.content.findMany({
      include: {
        textChunks: true,
      },
      orderBy: {
        [contentQueryDto.sortBy]: contentQueryDto.sortDirection,
      },
      skip: contentQueryDto.offset,
      take: contentQueryDto.limit,
      where: {
        createdAt: {
          gte: contentQueryDto.startDate,
          lte: contentQueryDto.endDate,
        },
        orgname,
        ...(contentQueryDto.searchTerm
          ? {
              OR: [{ name: { contains: contentQueryDto.searchTerm } }],
            }
          : undefined),
      },
    });
    return { count, results: content };
  }

  async findOne(id: string) {
    return this.prisma.content.findUniqueOrThrow({
      include: {
        textChunks: true,
      },
      where: { id },
    });
  }

  async incrementCredits(id: string, credits: number) {
    return this.prisma.content.update({
      data: {
        credits: {
          increment: credits,
        },
      },
      include: {
        textChunks: true,
      },
      where: {
        id,
      },
    });
  }

  async remove(orgname: string, id: string) {
    await this.prisma.content.delete({
      where: { id },
    });
  }

  async update(
    orgname: string,
    id: string,
    updateContentDto: UpdateContentDto
  ) {
    return this.prisma.content.update({
      data: {
        name: updateContentDto.name,
      },
      include: {
        textChunks: true,
      },
      where: {
        id,
      },
    });
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ContentUpdateInput) {
    return this.prisma.content.update({
      data: raw,
      include: {
        textChunks: true,
      },
      where: {
        id,
      },
    });
  }
}
