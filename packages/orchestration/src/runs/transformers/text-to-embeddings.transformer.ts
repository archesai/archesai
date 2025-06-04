// // import { ContentService } from "#content/content.service";
// import type { Logger } from '@archesai/core'
// // import { retry } from "#common/utils/retry";
// import type { ContentEntity } from '@archesai/domain'

// // import { OpenAiEmbeddingsService } from "#embeddings/embeddings.openai.service";

// // const chunkArray = <T>(array: T[], chunkSize: number): T[][] =>
// //   Array.from({ length: Math.ceil(array.length / chunkSize) }, (v, i) =>
// //     array.slice(i * chunkSize, i * chunkSize + chunkSize)
// //   );

// export const transformTextToEmbeddings = () // runId: string,
// // runInputContentIds: ContentEntity[],
// // logger: Logger,
// // contentService: ContentService,
// // openAiEmbeddingsService: OpenAiEmbeddingsService
// : ContentEntity[] => {
//   // const t1 = Date.now();
//   // let embeddings = [] satisfies {
//   //   embedding: number[];
//   //   tokens: number;
//   // }[];
//   // const textChunks = await contentService.findMany(content.id, {});
//   // const textContentChunks = chunkArray(textChunks.results, 100);
//   // for (const textContentChunk of textContentChunks) {
//   //   const embeddingsChunk = await retry(
//   //     logger,
//   //     async () =>
//   //       await openAiEmbeddingsService.createEmbeddings(
//   //         textContentChunk.map((x) => x.text)
//   //       ),
//   //     3
//   //   );
//   //   embeddings = embeddings.concat(embeddingsChunk);
//   // }

//   // logger.log(
//   //   `Created embeddings for ${content.name}.  Completed in ${
//   //     (Date.now() - t1) / 1000
//   //   }s`
//   // );

//   // const populatedTextContent = textChunks.results.map((textChunk, index) => {
//   //   return { ...textChunk, ...embeddings[index], textChunkId: textChunk.id };
//   // });

//   // const start = Date.now();
//   // await contentService.upsertVectors(
//   //   content.orgname,
//   //   content.id,
//   //   populatedTextContent
//   // );
//   // logger.log(
//   //   `Upserted embeddings for ${content.name}. Completed in ${
//   //     (Date.now() - start) / 1000
//   //   }s`
//   // );
//   return []
// }
