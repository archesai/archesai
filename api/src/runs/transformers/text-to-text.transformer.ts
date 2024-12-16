import { Logger } from '@nestjs/common'

import { retry } from '../../common/retry'
import { ContentService } from '../../content/content.service'
import { ContentEntity } from '../../content/entities/content.entity'
import { LLMService } from '../../llm/llm.service'
import { IToolRunProcess } from '../interfaces/tool-run-processor.interface'

export const transformTextToText: IToolRunProcess = async (
  runId: string,
  runInput: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  llmService: LLMService
): Promise<ContentEntity[]> => {
  logger.log(`Summarizing content for run ${runId}`)
  const start = Date.now()
  const c = runInput
    .map((x) => x.text)
    .filter((x) => x)
    .join(' ')
  logger.log(`Got first tokens for content for run ${runId}`)
  const { summary } = await retry(
    logger,
    async () => await llmService.createSummary(c),
    3
  )
  logger.log(`Got summary for content for run ${runId}`)

  logger.log('Summary saved. Completed in ' + (Date.now() - start) / 1000 + 's')

  logger.log(summary)

  const summaryContent = await contentService.create(runInput[0].orgname, {
    name: 'Summary Tool - ' + runInput.map((x) => x.name).join(', '),
    text: summary,
    labels: []
  })

  return [new ContentEntity(summaryContent)]
}
