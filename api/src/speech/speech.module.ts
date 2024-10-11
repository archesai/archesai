import { Module } from "@nestjs/common";

import { SpeechService } from "./speech.service";

@Module({
  exports: [SpeechService],
  providers: [SpeechService],
})
export class SpeechModule {}
