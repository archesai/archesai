// run.repository.ts
import { Injectable } from "@nestjs/common";
import { Run, RunStatus } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { RunQueryDto } from "./dto/run-query.dto";

@Injectable()
export class RunRepository
  implements BaseRepository<Run, undefined, RunQueryDto, undefined>
{
  constructor(private readonly prisma: PrismaService) {}

  async findAll(orgname: string, runQueryDto: RunQueryDto) {
    const count = await this.prisma.run.count({
      where: {
        createdAt: {
          gte: runQueryDto.startDate,
          lte: runQueryDto.endDate,
        },
        orgname,
      },
    });
    const results = await this.prisma.run.findMany({
      orderBy: {
        [runQueryDto.sortBy]: runQueryDto.sortDirection,
      },
      skip: runQueryDto.offset,
      take: runQueryDto.limit,
      where: {
        createdAt: {
          gte: runQueryDto.startDate,
          lte: runQueryDto.endDate,
        },
        orgname,
      },
    });
    return { count, results };
  }

  async findOne(orgname: string, id: string) {
    return this.prisma.run.findUniqueOrThrow({
      where: { id },
    });
  }

  async setCompletedAt(id: string, completedAt: Date) {
    return this.prisma.run.update({
      data: { completedAt },
      where: { id },
    });
  }

  async setProgress(id: string, progress: number) {
    return this.prisma.run.update({
      data: { progress },
      where: { id },
    });
  }

  async setRunError(id: string, error: string) {
    return this.prisma.run.update({
      data: { error },
      where: { id },
    });
  }

  async setStartedAt(id: string, startedAt: Date) {
    return this.prisma.run.update({
      data: { startedAt },
      where: { id },
    });
  }

  async updateStatus(id: string, status: RunStatus) {
    return this.prisma.run.update({
      data: { status },
      where: { id },
    });
  }
}
