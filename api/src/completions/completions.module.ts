import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { OpenAiCompletionsService } from "./completions.openai.service";

@Module({
  exports: [OpenAiCompletionsService],
  imports: [ConfigModule, HttpModule],
  providers: [OpenAiCompletionsService],
})
export class CompletionsModule {}
