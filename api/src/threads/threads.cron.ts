import { Injectable, Logger } from "@nestjs/common";
import { Cron, CronExpression } from "@nestjs/schedule";

import { ThreadsService } from "./threads.service"; // assuming you have a ThreadsService

@Injectable()
export class ThreadsCron {
  logger: Logger = new Logger(ThreadsCron.name);
  constructor(private readonly threadsService: ThreadsService) {}

  @Cron(CronExpression.EVERY_30_MINUTES)
  async cleanupUnused() {
    this.logger.log("Running cleanupUnused cronjob");
    const numRemoved = await this.threadsService.cleanupUnused();
    this.logger.log("Removed " + numRemoved + " unused threads");
  }
}
