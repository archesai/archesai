import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { Thread } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadEntity } from "./entities/thread.entity";
import { ThreadRepository } from "./thread.repository";

@Injectable()
export class ThreadsService extends BaseService<
  ThreadEntity,
  CreateThreadDto,
  undefined,
  ThreadRepository,
  {
    _count: {
      messages: number;
    };
  } & Thread
> {
  private readonly logger: Logger = new Logger("Threads Service");

  constructor(
    private threadRepository: ThreadRepository,
    private websocketsService: WebsocketsService
  ) {
    super(threadRepository);
  }

  async cleanupUnused() {
    return this.threadRepository.cleanupUnused();
  }

  async incrementCredits(orgname: string, threadId: string, credits: number) {
    return this.threadRepository.incrementCredits(orgname, threadId, credits);
  }

  protected toEntity(
    model: {
      _count: {
        messages: number;
      };
    } & Thread
  ): ThreadEntity {
    return new ThreadEntity(model);
  }

  async updateThreadName(orgname: string, threadId: string, name: string) {
    return this.threadRepository.updateThreadName(orgname, threadId, name);
  }
}
