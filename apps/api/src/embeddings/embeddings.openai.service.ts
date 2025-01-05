import { Injectable, Logger } from '@nestjs/common'
import OpenAI from 'openai'

import { EmbeddingsService } from './embeddings.service'
import { ArchesConfigService } from '../config/config.service'

@Injectable()
export class OpenAiEmbeddingsService implements EmbeddingsService {
  public openai: OpenAI
  private logger = new Logger(OpenAiEmbeddingsService.name)

  constructor(private configService: ArchesConfigService) {
    this.openai = new OpenAI({
      apiKey: this.configService.get('llm.token'),
      baseURL: this.configService.get('llm.endpoint'),
      organization: 'org-uCtGHWe8lpVBqo5thoryOqcS'
    })
  }

  async createEmbeddings(texts: string[]) {
    const start = Date.now()
    const { data, usage } = await this.openai.embeddings.create({
      input: texts,
      model:
        this.configService.get('llm.type') == 'openai'
          ? 'text-embedding-ada-002'
          : 'mxbai-embed-large'
    })
    const response = data.map((d) => {
      return {
        embedding: d.embedding,
        tokens: Math.ceil(usage.total_tokens / texts.length)
      }
    })
    this.logger.debug(
      `Embedded ${texts.length} texts with ${usage.total_tokens} in ${(Date.now() - start) / 1000}s`
    )
    return response
  }
}
