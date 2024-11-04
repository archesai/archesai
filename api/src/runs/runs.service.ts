// runs.service.ts
import { Injectable } from "@nestjs/common";
import { RunStatus } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { WebsocketsService } from "../websockets/websockets.service";
import { RunQueryDto } from "./dto/run-query.dto";
import { RunEntity } from "./entities/run.entity";
import { RunRepository } from "./run.repository";

@Injectable()
export class RunsService
  implements BaseService<RunEntity, undefined, RunQueryDto, undefined>
{
  constructor(
    private readonly runRepository: RunRepository,
    private websocketsService: WebsocketsService
  ) {}

  async findAll(orgname: string, runQueryDto: RunQueryDto) {
    const { count, results } = await this.runRepository.findAll(
      orgname,
      runQueryDto
    );
    const runEntities = results.map((run) => new RunEntity(run));
    return new PaginatedDto<RunEntity>({
      metadata: {
        limit: runQueryDto.limit,
        offset: runQueryDto.offset,
        totalResults: count,
      },
      results: runEntities,
    });
  }

  async findOne(orgname: string, id: string) {
    return new RunEntity(await this.runRepository.findOne(orgname, id));
  }

  async setProgress(id: string, progress: number) {
    const run = new RunEntity(
      await this.runRepository.setProgress(id, progress)
    );
    this.websocketsService.socket
      .to(run.orgname)
      .emit("update_progress", { ...run, orgname: run.orgname });
    return run;
  }

  async setRunError(id: string, error: string) {
    const run = new RunEntity(await this.runRepository.setRunError(id, error));
    this.websocketsService.socket.to(run.orgname).emit("update");
    return run;
  }

  async updateStatus(id: string, status: RunStatus) {
    switch (status) {
      case "COMPLETE":
        await this.runRepository.setCompletedAt(id, new Date());
        await this.runRepository.setProgress(id, 1);
        break;
      case "ERROR":
        await this.runRepository.setCompletedAt(id, new Date());
        break;
      case "PROCESSING":
        await this.runRepository.setStartedAt(id, new Date());
        break;
    }
    const run = new RunEntity(
      await this.runRepository.updateStatus(id, status)
    );
    this.websocketsService.socket.to(run.orgname).emit("update");
    return run;
  }
}
