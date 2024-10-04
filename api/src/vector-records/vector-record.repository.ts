import { Injectable } from "@nestjs/common";
import { VectorRecord } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { VectorRecordQueryDto } from "./dto/vector-record-query.dto";

@Injectable()
export class VectorRecordRepository
  implements
    BaseRepository<VectorRecord, undefined, VectorRecordQueryDto, undefined>
{
  constructor(private prisma: PrismaService) {}

  async findAll(
    orgname: string,
    vectorRecordQueryDto: VectorRecordQueryDto,
    contentId?: string
  ) {
    const count = await this.prisma.vectorRecord.count({
      where: {
        contentId,
        orgname,
        ...(vectorRecordQueryDto.searchTerm
          ? {
              OR: [{ text: { contains: vectorRecordQueryDto.searchTerm } }],
            }
          : undefined),
      },
    });
    const vectorRecord = await this.prisma.vectorRecord.findMany({
      orderBy: {
        [vectorRecordQueryDto.sortBy || "createdAt"]:
          vectorRecordQueryDto.sortDirection,
      },
      skip: vectorRecordQueryDto.offset,
      take: vectorRecordQueryDto.limit,
      where: {
        contentId,
        orgname,
        ...(vectorRecordQueryDto.searchTerm
          ? {
              OR: [{ text: { contains: vectorRecordQueryDto.searchTerm } }],
            }
          : undefined),
      },
    });
    return { count, results: vectorRecord };
  }

  async findOne(id: string) {
    return this.prisma.vectorRecord.findUniqueOrThrow({
      where: { id },
    });
  }
}
