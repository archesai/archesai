import { Injectable, Logger } from "@nestjs/common";
import { RunStatus } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { ContentEntity } from "../content/entities/content.entity";
import {
  TransformationEntity,
  TransformationModel,
} from "./entities/transformation.entity";
import { TransformationRepository } from "./transformation.repository";

@Injectable()
export class TransformationsService extends BaseService<
  TransformationEntity,
  any,
  any,
  TransformationRepository,
  TransformationModel
> {
  private logger = new Logger(TransformationsService.name);

  constructor(private transformationRepository: TransformationRepository) {
    super(transformationRepository);
  }

  async setOutputContent(transformationId: string, content: ContentEntity[]) {
    return this.toEntity(
      await this.transformationRepository.setOutputContent(
        transformationId,
        content
      )
    );
  }

  async setProgress(id: string, progress: number) {
    return this.toEntity(
      await this.transformationRepository.updateRaw(null, id, {
        progress,
      })
    );
  }

  async setRunError(id: string, error: string) {
    return this.toEntity(
      await this.transformationRepository.updateRaw(null, id, {
        error,
      })
    );
  }

  async setStatus(id: string, status: RunStatus) {
    switch (status) {
      case "COMPLETE":
        await this.transformationRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        await this.transformationRepository.updateRaw(null, id, {
          progress: 1,
        });
        break;
      case "ERROR":
        await this.transformationRepository.updateRaw(null, id, {
          completedAt: new Date(),
        });
        break;
      case "PROCESSING":
        await this.transformationRepository.updateRaw(null, id, {
          startedAt: new Date(),
        });
        break;
    }
    const run = await this.transformationRepository.updateRaw(null, id, {
      status,
    });
    return this.toEntity(run);
  }

  protected toEntity(model: TransformationModel): TransformationEntity {
    return new TransformationEntity(model);
  }
}
