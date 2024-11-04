import { TextChunksService } from "@/src/text-chunks/text-chunks.service";
import { Logger } from "@nestjs/common";

import { retry } from "../../common/retry";
import { ContentEntity } from "../../content/entities/content.entity";
import { OpenAiEmbeddingsService } from "../../embeddings/embeddings.openai.service";

const chunkArray = <T>(array: T[], chunkSize: number): T[][] =>
  Array.from({ length: Math.ceil(array.length / chunkSize) }, (v, i) =>
    array.slice(i * chunkSize, i * chunkSize + chunkSize)
  );

export const processCreateEmbeddings = async (
  content: ContentEntity,
  logger: Logger,
  textChunksService: TextChunksService,
  openAiEmbeddingsService: OpenAiEmbeddingsService
) => {
  const t1 = Date.now();
  let embeddings = [] as {
    embedding: number[];
    tokens: number;
  }[];
  const textChunks = await textChunksService.findAll(content.id, {});
  const textContentChunks = chunkArray(textChunks.results, 100);
  for (const textContentChunk of textContentChunks) {
    const embeddingsChunk = await retry(
      logger,
      async () =>
        await openAiEmbeddingsService.createEmbeddings(
          textContentChunk.map((x) => x.text)
        ),
      3
    );
    embeddings = embeddings.concat(embeddingsChunk);
  }

  logger.log(
    `Created embeddings for ${content.name}.  Completed in ${
      (Date.now() - t1) / 1000
    }s`
  );

  const populatedTextContent = textChunks.results.map((textChunk, index) => {
    return { ...textChunk, ...embeddings[index], textChunkId: textChunk.id };
  });

  const start = Date.now();
  await textChunksService.upsertVectors(
    content.orgname,
    content.id,
    populatedTextContent
  );
  logger.log(
    `Upserted embeddings for ${content.name}. Completed in ${
      (Date.now() - start) / 1000
    }s`
  );
};
