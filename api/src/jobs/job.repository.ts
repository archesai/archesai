import { Injectable } from "@nestjs/common";
import { JobStatus } from "@prisma/client";
import { Job } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { JobQueryDto } from "./dto/job-query.dto";

@Injectable()
export class JobRepository
  implements BaseRepository<Job, undefined, JobQueryDto, undefined>
{
  constructor(private readonly prisma: PrismaService) {}

  async findAll(orgname: string, jobQueryDto: JobQueryDto) {
    const count = await this.prisma.job.count({
      where: { ...jobQueryDto, orgname },
    });
    const results = await this.prisma.job.findMany({
      orderBy: {
        [jobQueryDto.sortBy]: jobQueryDto.sortDirection,
      },
      skip: jobQueryDto.offset,
      take: jobQueryDto.limit,
      where: { ...jobQueryDto, orgname },
    });
    return { count, results };
  }

  findOne(orgname: string, id: string) {
    return this.prisma.job.findUnique({
      where: { id },
    });
  }

  async remove(orgname: string, id: string) {
    await this.prisma.job.delete({
      where: { id },
    });
  }

  async setCompletedAt(id: string, completedAt: Date) {
    return this.prisma.job.update({
      data: { completedAt },
      where: { id },
    });
  }

  async setProgress(id: string, progress: number) {
    return this.prisma.job.update({
      data: { progress },
      where: { id },
    });
  }

  async setStartedAt(id: string, startedAt: Date) {
    return this.prisma.job.update({
      data: { startedAt },
      where: { id },
    });
  }

  async updateStatus(id: string, status: JobStatus) {
    return this.prisma.job.update({
      data: { status },
      where: { id },
    });
  }
}
