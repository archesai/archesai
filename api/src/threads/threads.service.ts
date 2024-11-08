import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadEntity, ThreadModelWithCount } from "./entities/thread.entity";
import { ThreadRepository } from "./thread.repository";

@Injectable()
export class ThreadsService extends BaseService<
  ThreadEntity,
  CreateThreadDto,
  undefined,
  ThreadRepository,
  ThreadModelWithCount
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

  protected toEntity(model: ThreadModelWithCount): ThreadEntity {
    return new ThreadEntity(model);
  }
}
