// import path from 'node:path'
// import type { Logger } from '@archesai/core'
// import type { RunpodService } from '@archesai/intelligence'
// import type { ArtifactEntity } from '@archesai/schemas'
// import type { StorageService } from '@archesai/storage'

// import type { ArtifactsService } from '#artifacts/artifacts.service'
// import type { Transformer } from '#runs/types/transformer'

// export const transformTextToImage: Transformer = async (
//   runId: string,
//   inputs: ArtifactEntity[],
//   logger: Logger,
//   artifactsService: ArtifactsService,
//   runpodService: RunpodService,
//   storageService: StorageService
// ): Promise<ArtifactEntity[]> => {
//   logger.log(
//     {
//       inputs,
//       runId
//     },
//     `processing text to image`
//   )
//   const input = {
//     input: {
//       prompt: 'a man running in a circle'
//     }
//   }
//   if (!inputs[0]) {
//     throw new Error('not enough inputs')
//   }

//   const image_url = await runpodService.run(runId, 'y55cw5fvbum8q6', input)

//   const base64String = image_url.replace(/^data:image\/\w+;base64,/, '')

//   // Convert the remaining base64 string to a buffer
//   const buffer = Buffer.from(base64String, 'base64')
//   const bucketPath = `images/${runId}.png`

//   const fileEntity = await storageService.uploadFromFile(path, {
//     buffer: buffer,
//     mimetype: 'image/png',
//     originalname: path.basename(bucketPath)
//   })
//   logger.log({ fileEntity, runId }, `text-to-image complete`)

//   const content = await artifactsService.create({
//     credits: 0,
//     description: 'Text to Speech Tool -' + inputs.map((x) => x.name).join(', '),
//     labels: [],
//     mimeType: 'image/png',
//     name: 'Text to Speech Tool -' + inputs.map((x) => x.name).join(', '),
//     organizationId: inputs[0].organizationId,
//     text: null,
//     url: fileEntity.read ?? null
//   })

//   return [content]
// }
