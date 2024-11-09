import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { GranularCount, GranularSum } from "../common/dto/aggregated-field.dto";
import { PrismaService } from "../prisma/prisma.service";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadAggregates } from "./dto/thread-aggregates.dto";
import { ThreadQueryDto } from "./dto/thread-query.dto";
import { ThreadModel } from "./entities/thread.entity";

@Injectable()
export class ThreadRepository extends BaseRepository<
  ThreadModel,
  CreateThreadDto,
  undefined,
  Prisma.ThreadInclude,
  Prisma.ThreadSelect,
  Prisma.ThreadUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.thread);
  }

  async findAll(orgname: string, threadQueryDto: ThreadQueryDto) {
    const whereConditions = {
      createdAt: {
        gte: threadQueryDto.startDate,
        lte: threadQueryDto.endDate,
      },
      orgname,
    };
    if (threadQueryDto.filters) {
      threadQueryDto.filters.forEach((filter) => {
        whereConditions[filter.field] = { [filter.operator]: filter.value };
      });
    }

    const count = await this.prisma.thread.count({
      where: whereConditions,
    });

    let aggregates: ThreadAggregates = null as {
      credits: GranularSum[];
      threadsCreated: GranularCount[];
    };

    if (threadQueryDto.aggregates) {
      const rawAggregates = await this.prisma.$queryRaw`
          SELECT 
              DATE_TRUNC(${threadQueryDto.aggregateGranularity}, "createdAt") AS day, 
              COUNT(*) AS count,
              COALESCE(SUM("credits"), 0)::numeric AS "dailyCredits"
          FROM "Thread"
          WHERE "orgname" = ${orgname}
          GROUP BY day
          ORDER BY day;
      `;

      // Convert the daily data to the desired format

      aggregates = {
        credits: (rawAggregates as any).map((record) => ({
          from: record.day,
          sum: Number(record.dailyCredits),
          to: new Date(
            new Date(record.day).getTime() + 24 * 60 * 60 * 1000 - 1
          ), // end of the day
        })),
        threadsCreated: (rawAggregates as any).map((record) => ({
          count: Number(record.count),
          from: record.day,
          to: new Date(
            new Date(record.day).getTime() + 24 * 60 * 60 * 1000 - 1
          ), // end of the day
        })),
      };
    }

    const results = await this.prisma.thread.findMany({
      orderBy: {
        [threadQueryDto.sortBy]: threadQueryDto.sortDirection,
      },
      skip: threadQueryDto.offset,
      take: threadQueryDto.limit,
      where: whereConditions,
    });
    return { aggregates, count, results };
  }
}
