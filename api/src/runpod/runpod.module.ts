import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";

import { TransformationsModule } from "../transformations/transformations.module";
import { RunpodService } from "./runpod.service";

@Module({
  exports: [RunpodService],
  imports: [TransformationsModule, HttpModule],
  providers: [RunpodService],
})
export class RunpodModule {}
