import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";

import { LLMService } from "./llm.service";

@Module({
  exports: [LLMService],
  imports: [HttpModule],
  providers: [LLMService],
})
export class LLMModule {}
