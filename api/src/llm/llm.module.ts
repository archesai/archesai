import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { LLMService } from "./llm.service";

@Module({
  exports: [LLMService],
  imports: [ConfigModule, HttpModule],
  providers: [LLMService],
})
export class LLMModule {}
