import { Injectable } from "@nestjs/common";
import { Job, JobStatus } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { JobQueryDto } from "./dto/job-query.dto";
import { JobEntity } from "./entities/job.entity";
import { JobRepository } from "./job.repository";

@Injectable()
export class JobsService
  implements BaseService<Job, undefined, JobQueryDto, undefined>
{
  constructor(
    private readonly jobRepository: JobRepository,
    private websocketsService: WebsocketsService
  ) {}
  async findAll(orgname: string, jobQueryDto: JobQueryDto) {
    return this.jobRepository.findAll(orgname, jobQueryDto);
  }

  async findOne(orgname: string, id: string) {
    return new JobEntity(await this.jobRepository.findOne(orgname, id));
  }

  async remove(orgname: string, id: string) {
    this.websocketsService.socket.to(orgname).emit("update");
    await this.jobRepository.remove(orgname, id);
  }

  async setJobError(id: string, error: string) {
    const job = new JobEntity(await this.jobRepository.setJobError(id, error));
    this.websocketsService.socket.to(job.orgname).emit("update");
    return job;
  }

  async setProgress(id: string, progress: number) {
    const job = new JobEntity(
      await this.jobRepository.setProgress(id, progress)
    );
    this.websocketsService.socket
      .to(job.orgname)
      .emit("update_progress", { ...job, orgname: job.orgname });
    return job;
  }

  async updateStatus(id: string, status: JobStatus) {
    switch (status) {
      case "COMPLETE":
        await this.jobRepository.setCompletedAt(id, new Date());
        await this.jobRepository.setProgress(id, 1);
        break;
      case "ERROR":
        await this.jobRepository.setCompletedAt(id, new Date());
        break;
      case "PROCESSING":
        await this.jobRepository.setStartedAt(id, new Date());
        break;
    }
    const job = new JobEntity(
      await this.jobRepository.updateStatus(id, status)
    );
    this.websocketsService.socket.to(job.orgname).emit("update");
    return job;
  }
}
