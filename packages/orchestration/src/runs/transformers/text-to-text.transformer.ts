// import type { Logger } from '@archesai/core'
// import type { LlmService } from '@archesai/intelligence'
// import type { ArtifactEntity } from '@archesai/schemas'
// import { retry } from '@archesai/core'
// import { contentEntitySchema } from '@archesai/schemas'

// import type { ArtifactsService } from '#artifacts/artifacts.service'
// import type { Transformer } from '#runs/types/transformer'

// export const transformTextToText: Transformer = async (
//   runId: string,
//   runInput: ArtifactEntity[],
//   logger: Logger,
//   artifactsService: ArtifactsService,
//   llmService: LlmService
// ): Promise<ArtifactEntity[]> => {
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

//   const summaryContent = await artifactsService.create({
//     credits: 0,
//     description: '',
//     labels: [],
//     mimeType: 'text/plain',
//     name: 'Summary Tool - ' + runInput.map((x) => x.name).join(', '),
//     organizationId: runInput[0].organizationId,
//     text: summary,
//     url: null
//   })

//   return [contentEntitySchema.parse(summaryContent)]
// }
