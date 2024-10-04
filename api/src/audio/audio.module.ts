import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";

import { StorageModule } from "../storage/storage.module";
import { AudioService } from "./audio.service";
import { KeyframesService } from "./keyframes.service";

@Module({
  exports: [AudioService],
  imports: [StorageModule.forRoot(), HttpModule],
  providers: [AudioService, KeyframesService],
})
export class AudioModule {}
