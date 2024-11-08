// runs.service.ts
import { InjectFlowProducer, InjectQueue } from "@nestjs/bullmq";
import { BadRequestException, Injectable } from "@nestjs/common";
import { Run, RunStatus } from "@prisma/client";
import { FlowProducer, Queue } from "bullmq";

import { BaseService } from "../common/base.service";
import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { RunToolDto } from "../tools/dto/run-tool.dto";
import { ToolEntity } from "../tools/entities/tool.entity";
import { WebsocketsService } from "../websockets/websockets.service";
import { RunEntity } from "./entities/run.entity";
import { RunDetailedEntity } from "./entities/run-detailed.entity";
import { RunRepository } from "./run.repository";

@Injectable()
export class RunsService extends BaseService<
  RunEntity,
  undefined,
  undefined,
  RunRepository,
  Run
> {
  constructor(
    private readonly runRepository: RunRepository,
    private websocketsService: WebsocketsService,
    @InjectFlowProducer("flow") private readonly flowProducer: FlowProducer,
    @InjectQueue("run") private readonly runQueue: Queue,
    private contentService: ContentService
  ) {
    super(runRepository);
  }

  async addRunInputContent(id: string, contents: ContentEntity[]) {
    const run = new RunEntity(
      await this.runRepository.addOutputContent(id, contents)
    );
    this.websocketsService.socket.to(run.orgname).emit("update", {
      queryKey: ["organizations", run.orgname, "runs"],
    });
    return run;
  }

  async findOne(orgname: string, id: string): Promise<RunDetailedEntity> {
    return new RunDetailedEntity(await this.runRepository.findOne(orgname, id));
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

  async runTool(orgname: string, tool: ToolEntity, runToolDto: RunToolDto) {
    if (!!runToolDto.runInputContentIds?.length) {
      // vefify that the content exists
      for (const contentId of runToolDto.runInputContentIds) {
        await this.contentService.findOne(orgname, contentId);
      }
    } else if (runToolDto.text) {
      const content = await this.contentService.create(orgname, {
        name: "Tool Input Text - " + tool.id,
        text: runToolDto.text,
      });
      runToolDto.runInputContentIds = [content.id];
    } else if (runToolDto.url) {
      const content = await this.contentService.create(orgname, {
        name: "Tool Input URL - " + tool.id,
        url: runToolDto.url,
      });
      runToolDto.runInputContentIds = [content.id];
    } else {
      throw new BadRequestException("No input provided");
    }

    const run = await this.runRepository.createToolRun(
      orgname,
      tool.id,
      runToolDto
    );

    const runInputContents: ContentEntity[] = [];
    for (const contentId of runToolDto.runInputContentIds) {
      runInputContents.push(
        await this.contentService.findOne(orgname, contentId)
      );
    }

    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "runs"],
    });

    await this.runQueue.add(
      tool.toolBase,
      {
        runInputContents,
      },
      {
        jobId: run.id,
      }
    );

    return new RunEntity(run);
  }

  async setProgress(id: string, progress: number) {
    const run = new RunEntity(
      await this.runRepository.updateRaw(null, id, {
        progress,
      })
    );
    this.websocketsService.socket
      .to(run.orgname)
      .emit("update_progress", { ...run, orgname: run.orgname });
    return run;
  }

  async setRunError(id: string, error: string) {
    const run = new RunEntity(
      await this.runRepository.updateRaw(null, id, {
        error,
      })
    );
    this.websocketsService.socket.to(run.orgname).emit("update", {
      queryKey: ["organizations", run.orgname, "runs"],
    });
    return run;
  }

  async setStatus(id: string, status: RunStatus) {
    switch (status) {
      case "COMPLETE":
        await this.runRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        await this.runRepository.updateRaw(null, id, {
          progress: 1,
        });
        break;
      case "ERROR":
        await this.runRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        break;
      case "PROCESSING":
        await this.runRepository.updateRaw(null, id, {
          startedAt: new Date(),
        });
        break;
    }
    const run = await this.runRepository.updateRaw(null, id, {
      status,
    });
    this.websocketsService.socket.to(run.orgname).emit("update", {
      queryKey: ["organizations", run.orgname, "runs"],
    });
    return this.toEntity(run);
  }

  protected toEntity(model: Run): RunEntity {
    return new RunEntity(model);
  }
}
