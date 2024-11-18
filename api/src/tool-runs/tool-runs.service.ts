import { Injectable, Logger } from "@nestjs/common";
import { RunStatus } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { CreateRunDto } from "../common/dto/create-run.dto";
import { ContentEntity } from "../content/entities/content.entity";
import { WebsocketsService } from "../websockets/websockets.service";
import { ToolRunEntity, ToolRunModel } from "./entities/tool-run.entity";
import { ToolRunRepository } from "./tool-run.repository";

@Injectable()
export class ToolRunsService extends BaseService<
  ToolRunEntity,
  CreateRunDto,
  any,
  ToolRunRepository,
  ToolRunModel
> {
  private logger = new Logger(ToolRunsService.name);

  constructor(
    private toolRunRepository: ToolRunRepository,
    private websocketsService: WebsocketsService
  ) {
    super(toolRunRepository);
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "tool-runs"],
    });
  }

  async setOutputContent(toolRunId: string, content: ContentEntity[]) {
    return this.toEntity(
      await this.toolRunRepository.setOutputContent(toolRunId, content)
    );
  }

  async setProgress(id: string, progress: number) {
    return this.toEntity(
      await this.toolRunRepository.updateRaw(null, id, {
        progress,
      })
    );
  }

  async setRunError(id: string, error: string) {
    return this.toEntity(
      await this.toolRunRepository.updateRaw(null, id, {
        error,
      })
    );
  }

  async setStatus(id: string, status: RunStatus) {
    switch (status) {
      case "COMPLETE":
        await this.toolRunRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        await this.toolRunRepository.updateRaw(null, id, {
          progress: 1,
        });
        break;
      case "ERROR":
        await this.toolRunRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        break;
      case "PROCESSING":
        await this.toolRunRepository.updateRaw(null, id, {
          startedAt: new Date(),
        });
        break;
    }
    const run = await this.toolRunRepository.updateRaw(null, id, {
      status,
    });
    return this.toEntity(run);
  }

  protected toEntity(model: ToolRunModel): ToolRunEntity {
    return new ToolRunEntity(model);
  }
}
