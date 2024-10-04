import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { OpenAiEmbeddingsService } from "./embeddings.openai.service";

@Module({
  exports: [OpenAiEmbeddingsService],
  imports: [ConfigModule, HttpModule],
  providers: [OpenAiEmbeddingsService],
})
export class EmbeddingsModule {}
