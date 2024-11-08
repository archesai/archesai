import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { Thread } from "@prisma/client";

import { BaseService } from "../common/base.service";
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

  constructor(private threadRepository: ThreadRepository) {
    super(threadRepository);
  }

  async cleanupUnused() {
    return this.threadRepository.cleanupUnused();
  }

  async setTitle(orgname: string, id: string, title: string) {
    return this.toEntity(
      await this.threadRepository.updateRaw(orgname, id, {
        name: title,
      })
    );
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
}
