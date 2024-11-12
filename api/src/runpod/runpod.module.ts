import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";

import { ToolRunsModule } from "../tool-runs/tool-runs.module";
import { RunpodService } from "./runpod.service";

@Module({
  exports: [RunpodService],
  imports: [ToolRunsModule, HttpModule],
  providers: [RunpodService],
})
export class RunpodModule {}
