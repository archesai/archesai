import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { GranularCount, GranularSum } from "../common/dto/aggregated-field.dto";
import { PrismaService } from "../prisma/prisma.service";
import { CreateLabelDto } from "./dto/create-label.dto";
import { LabelAggregates } from "./dto/label-aggregates.dto";
import { LabelQueryDto } from "./dto/label-query.dto";
import { LabelModel } from "./entities/label.entity";

@Injectable()
export class LabelRepository extends BaseRepository<
  LabelModel,
  CreateLabelDto,
  undefined,
  Prisma.LabelInclude,
  Prisma.LabelSelect,
  Prisma.LabelUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.label);
  }

  async findAll(orgname: string, labelQueryDto: LabelQueryDto) {
    const whereConditions = {
      createdAt: {
        gte: labelQueryDto.startDate,
        lte: labelQueryDto.endDate,
      },
      orgname,
    };
    if (labelQueryDto.filters) {
      labelQueryDto.filters.forEach((filter) => {
        whereConditions[filter.field] = { [filter.operator]: filter.value };
      });
    }

    const count = await this.prisma.label.count({
      where: whereConditions,
    });

    let aggregates: LabelAggregates = null as {
      credits: GranularSum[];
      labelsCreated: GranularCount[];
    };

    if (labelQueryDto.aggregates) {
      const rawAggregates = await this.prisma.$queryRaw`
          SELECT 
              DATE_TRUNC(${labelQueryDto.aggregateGranularity}, "createdAt") AS day, 
              COUNT(*) AS count,
              COALESCE(SUM("credits"), 0)::numeric AS "dailyCredits"
          FROM "Label"
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
        labelsCreated: (rawAggregates as any).map((record) => ({
          count: Number(record.count),
          from: record.day,
          to: new Date(
            new Date(record.day).getTime() + 24 * 60 * 60 * 1000 - 1
          ), // end of the day
        })),
      };
    }

    const results = await this.prisma.label.findMany({
      orderBy: {
        [labelQueryDto.sortBy]: labelQueryDto.sortDirection,
      },
      skip: labelQueryDto.offset,
      take: labelQueryDto.limit,
      where: whereConditions,
    });
    return { aggregates, count, results };
  }
}
