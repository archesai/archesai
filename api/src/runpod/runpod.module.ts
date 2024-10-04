import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";

import { JobsModule } from "../jobs/jobs.module";
import { RunpodService } from "./runpod.service";

@Module({
  exports: [RunpodService],
  imports: [JobsModule, HttpModule],
  providers: [RunpodService],
})
export class RunpodModule {}
