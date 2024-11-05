// runs.service.ts
import { InjectFlowProducer, InjectQueue } from "@nestjs/bullmq";
import { Injectable } from "@nestjs/common";
import { RunStatus } from "@prisma/client";
import { FlowProducer, Queue } from "bullmq";

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
    private websocketsService: WebsocketsService,
    @InjectFlowProducer("flow") private readonly flowProducer: FlowProducer,
    @InjectQueue("run") private readonly runQueue: Queue
  ) {}

  async createPipelineRun(
    orgname: string,
    pipelineId: string,
    runInputContentIds: string[]
  ) {
    const run = await this.runRepository.createPipelineRun(
      orgname,
      pipelineId,
      runInputContentIds
    );
    this.websocketsService.socket.to(orgname).emit("update");

    await this.flowProducer.add({
      data: {
        // content: content,
        toolId: "extract-text",
      },
      name: "extract-text",
      queueName: "tool",
    });
    return new RunEntity(run);
  }

  async createToolRun(
    orgname: string,
    toolId: string,
    runInputContentIds: string[]
  ) {
    const run = await this.runRepository.createToolRun(
      orgname,
      toolId,
      runInputContentIds
    );
    this.websocketsService.socket.to(orgname).emit("update");
    const runEntity = new RunEntity(run);
    // await this.runQueue.add("as", runEntity);
    return runEntity;
  }

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
