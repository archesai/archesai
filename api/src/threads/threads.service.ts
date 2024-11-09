import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { ContentService } from "../content/content.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadEntity, ThreadModel } from "./entities/thread.entity";
import { ThreadRepository } from "./thread.repository";

@Injectable()
export class ThreadsService extends BaseService<
  ThreadEntity,
  CreateThreadDto,
  undefined,
  ThreadRepository,
  ThreadModel
> {
  private readonly logger: Logger = new Logger("Threads Service");

  constructor(
    private threadRepository: ThreadRepository,
    private contentService: ContentService,
    private websocketsService: WebsocketsService
  ) {
    super(threadRepository);
  }

  async setTitle(orgname: string, id: string, title: string) {
    return this.toEntity(
      await this.threadRepository.updateRaw(orgname, id, {
        name: title,
      })
    );
  }

  protected toEntity(model: ThreadModel): ThreadEntity {
    return new ThreadEntity(model);
  }
}
