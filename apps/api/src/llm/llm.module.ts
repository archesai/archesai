import { HttpModule } from '@nestjs/axios'
import { Module } from '@nestjs/common'

import { LlmService } from './llm.service'

@Module({
  exports: [LlmService],
  imports: [HttpModule],
  providers: [LlmService]
})
export class LlmModule {}
