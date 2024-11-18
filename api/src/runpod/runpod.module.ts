import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";

import { RunsModule } from "../runs/runs.module";
import { RunpodService } from "./runpod.service";

@Module({
  exports: [RunpodService],
  imports: [RunsModule, HttpModule],
  providers: [RunpodService],
})
export class RunpodModule {}
