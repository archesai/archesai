import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";

import { OpenAiEmbeddingsService } from "./embeddings.openai.service";

@Module({
  exports: [OpenAiEmbeddingsService],
  imports: [HttpModule],
  providers: [OpenAiEmbeddingsService],
})
export class EmbeddingsModule {}
