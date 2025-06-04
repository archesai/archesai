// import type { Logger } from '@archesai/core'
// import type { LlmService } from '@archesai/intelligence'
// import type { ContentEntity } from '@archesai/domain'
// import { retry } from '@archesai/core'
// import { contentEntitySchema } from '@archesai/domain'

// import type { ContentService } from '#content/content.service'
// import type { Transformer } from '#runs/types/transformer'

// export const transformTextToText: Transformer = async (
//   runId: string,
//   runInput: ContentEntity[],
//   logger: Logger,
//   contentService: ContentService,
//   llmService: LlmService
// ): Promise<ContentEntity[]> => {
//   logger.log(
//     {
//       runId
//     },
//     `summarizing content`
//   )
//   if (!runInput[0]?.text) {
//     throw new Error('no text provided')
//   }
//   const c = runInput
//     .map((x) => x.text)
//     .filter((x) => x)
//     .join(' ')
//   logger.log(
//     {
//       runId
//     },
//     `got first tokens for content`
//   )
//   const { summary } = await retry(
//     logger,
//     async () => await llmService.createSummary(c),
//     3
//   )
//   logger.log(
//     {
//       runId,
//       summary
//     },
//     `got summary for content`
//   )

//   const summaryContent = await contentService.create({
//     credits: 0,
//     description: '',
//     labels: [],
//     mimeType: 'text/plain',
//     name: 'Summary Tool - ' + runInput.map((x) => x.name).join(', '),
//     orgname: runInput[0].orgname,
//     text: summary,
//     url: null
//   })

//   return [contentEntitySchema.parse(summaryContent)]
// }
