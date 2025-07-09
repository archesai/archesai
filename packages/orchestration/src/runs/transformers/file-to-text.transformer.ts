// import type { Logger } from '@archesai/core'

// import { UnstructuredLoader } from '@langchain/community/document_loaders/fs/unstructured'

// import type { ConfigService } from '@archesai/core'
// import type { ArtifactEntity } from '@archesai/schemas'
// import type { StorageService } from '@archesai/storage'

// import type { ArtifactsService } from '#artifacts/artifacts.service'
// import type { Transformer } from '#runs/types/transformer'

// export const transformFileToText: Transformer = async (
//   runId: string,
//   inputs: ArtifactEntity[],
//   logger: Logger,
//   artifactsService: ArtifactsService,
//   configService: ConfigService,
//   storageService: StorageService
// ): Promise<ArtifactEntity[]> => {
//   logger.log(
//     {
//       inputs,
//       runId
//     },
//     `extracing text`
//   )

//   let content = inputs[0]
//   if (!content?.url) {
//     throw new Error('No url provided')
//   }
//   const buffer = await storageService.downloadToBuffer(content.url)
//   const loader = new UnstructuredLoader(
//     {
//       buffer: buffer,
//       fileName: content.url
//     },
//     {
//       apiUrl: configService.get('unstructured.endpoint')!
//     }
//   )

//   const docs = await loader.load()
//   const textContent = docs.map((doc) => {
//     return {
//       page: doc.metadata.page,
//       text: doc.pageContent,
//       tokens: doc.metadata.tokens
//     }
//   })

//   const title = 'UNIMPLEMENTED TITLE'

//   const sanitizedTextContent = textContent.map((data) => ({
//     ...data,
//     text: data.text
//       .replaceAll(/\0/g, '')
//       .replaceAll(/[^ -~\u00A0-\uD7FF\uE000-\uFDCF\uFDF0-\uFFFD\n]/g, '')
//   }))

//   // update name
//   if (!title.includes('http')) {
//     content = await artifactsService.update(content.id, {
//       name: title
//     })
//   }

//   const chunkedContent = await Promise.all(
//     sanitizedTextContent.map((data, i) =>
//       artifactsService.create({
//         credits: 0,
//         description: '',
//         labels: [],
//         mimeType: 'text/plain',
//         name: `${content.name} - Page ${data.page} - Index ${i}`,
//         orgname: content.orgname,
//         previewImage: '',
//         text: data.text,
//         url: null
//       })
//     )
//   )

//   return chunkedContent
// }
