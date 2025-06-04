// import type { Logger } from '@archesai/core'

// import type { SpeechService } from '@archesai/intelligence'
// import type { ArtifactEntity } from '@archesai/domain'
// import type { StorageService } from '@archesai/storage'
// import { contentEntitySchema } from '@archesai/domain'

// import type { ArtifactsService } from '#artifacts/artifacts.service'
// import type { Transformer } from '#runs/types/transformer'

// export const transformTextToSpeech: Transformer = async (
//   runId: string,
//   inputs: ArtifactEntity[],
//   logger: Logger,
//   artifactsService: ArtifactsService,
//   storageService: StorageService,
//   speechService: SpeechService
// ): Promise<ArtifactEntity[]> => {
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

//   const content = await artifactsService.create({
//     labels: [],
//     mimeType: 'audio/mpeg',
//     name: 'Text to Speech Tool -' + inputs.map((x) => x.name).join(', '),
//     orgname: inputs[0].orgname,
//     text: null,
//     url: fileEntity.read ?? null
//   })

//   return [contentEntitySchema.parse([content])]
// }
