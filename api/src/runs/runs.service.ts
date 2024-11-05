// runs.service.ts
import { InjectFlowProducer, InjectQueue } from "@nestjs/bullmq";
import { BadRequestException, Injectable } from "@nestjs/common";
import { RunStatus } from "@prisma/client";
import { FlowProducer, Queue } from "bullmq";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { ContentService } from "../content/content.service";
import { RunToolDto } from "../tools/dto/run-tool.dto";
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
    @InjectQueue("run") private readonly runQueue: Queue,
    private contentService: ContentService
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

  async runPipeline(
    orgname: string,
    pipelineId: string,
    runInputContentIds: string[]
  ) {
    const run = await this.runRepository.createPipelineRun(
      orgname,
      pipelineId,
      runInputContentIds
    );
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "piplines"],
    });

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

  async runTool(orgname: string, toolId: string, runToolDto: RunToolDto) {
    if (runToolDto.runInputContentIds) {
      // vefify that the content exists
      for (const contentId of runToolDto.runInputContentIds) {
        await this.contentService.findOne(contentId);
      }
    } else if (runToolDto.text) {
      const content = await this.contentService.create(orgname, {
        name: "Input Text",
        url: runToolDto.text,
      });
      runToolDto.runInputContentIds = [content.id];
    } else if (runToolDto.url) {
      const content = await this.contentService.create(orgname, {
        name: "Input URL",
        url: runToolDto.url,
      });
      runToolDto.runInputContentIds = [content.id];
    } else {
      throw new BadRequestException("No input provided");
    }

    const run = await this.runRepository.createToolRun(
      orgname,
      toolId,
      runToolDto
    );

    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "runs"],
    });
    const runEntity = new RunEntity(run);
    // await this.runQueue.add("as", runEntity);
    return runEntity;
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
