// import type { Logger } from '@archesai/core'

// import type { SpeechService } from '@archesai/ai'
// import type { ContentEntity } from '@archesai/domain'
// import type { StorageService } from '@archesai/storage'
// import { contentEntitySchema } from '@archesai/domain'

// import type { ContentService } from '#content/content.service'
// import type { Transformer } from '#runs/types/transformer'

// export const transformTextToSpeech: Transformer = async (
//   runId: string,
//   inputs: ContentEntity[],
//   logger: Logger,
//   contentService: ContentService,
//   storageService: StorageService,
//   speechService: SpeechService
// ): Promise<ContentEntity[]> => {
//   logger.log(
//     {
//       runId
//     },
//     `processing text to speech`
//   )
//   if (!inputs[0]?.text) {
//     throw new Error('no text provided')
//   }

//   const buffer = await speechService.generateSpeech(
//     inputs.map((x) => x.text).join(' ')
//   )

//   const fileEntity = await storageService.uploadFromFile(
//     `contents/${runId}.mp3`,
//     {
//       buffer,
//       mimetype: 'audio/mpeg',
//       originalname: `${runId}.mp3`
//     }
//   )

//   const content = await contentService.create({
//     labels: [],
//     mimeType: 'audio/mpeg',
//     name: 'Text to Speech Tool -' + inputs.map((x) => x.name).join(', '),
//     orgname: inputs[0].orgname,
//     text: null,
//     url: fileEntity.read ?? null
//   })

//   return [contentEntitySchema.parse([content])]
// }
