import { Injectable } from "@nestjs/common";
import { Prisma, Thread } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { GranularCount, GranularSum } from "../common/dto/aggregated-field.dto";
import { PrismaService } from "../prisma/prisma.service";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadAggregates } from "./dto/thread-aggregates.dto";
import { ThreadQueryDto } from "./dto/thread-query.dto";

@Injectable()
export class ThreadRepository extends BaseRepository<
  {
    _count: {
      messages: number;
    };
  } & Thread,
  CreateThreadDto,
  undefined,
  Prisma.ThreadInclude,
  Prisma.ThreadSelect,
  Prisma.ThreadUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.thread, {
      _count: {
        select: {
          messages: true,
        },
      },
    });
  }

  async cleanupUnused() {
    // First, fetch all threads with no messagess.
    const threads = await this.prisma.thread.findMany({
      where: {
        messages: {
          none: {},
        },
      },
    });

    // Then, delete each thread one by one.
    for (const thread of threads) {
      await this.prisma.thread.delete({
        where: {
          id: thread.id,
        },
      });
    }

    // Optionally, return the number of deleted threads.
    return threads.length;
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
      include: {
        _count: {
          select: {
            messages: true,
          },
        },
      },
      orderBy: {
        [threadQueryDto.sortBy]: threadQueryDto.sortDirection,
      },
      skip: threadQueryDto.offset,
      take: threadQueryDto.limit,
      where: whereConditions,
    });
    return { aggregates, count, results };
  }

  async incrementCredits(orgname: string, threadId: string, credits: number) {
    return this.prisma.thread.update({
      data: {
        credits: {
          increment: Math.ceil(credits),
        },
      },
      where: { id: threadId },
    });
  }

  async updateThreadName(orgname: string, threadId: string, name: string) {
    return this.prisma.thread.update({
      data: {
        name,
      },
      where: { id: threadId },
    });
  }
}
