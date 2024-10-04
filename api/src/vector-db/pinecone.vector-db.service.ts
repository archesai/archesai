import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { Index, Pinecone, RecordMetadata } from "@pinecone-database/pinecone";

import { BaseVectorDBService, VectorDBService } from "./vector-db.service";

@Injectable()
export class PineconeVectorDBService
  extends BaseVectorDBService
  implements VectorDBService
{
  index: Index<RecordMetadata>;
  pinecone: Pinecone;
  constructor(private configService: ConfigService) {
    super();
    this.pinecone = new Pinecone({
      apiKey: this.configService.get<string>("PINECONE_API_KEY"),
    });
    const index = this.configService.get<string>("PINECONE_INDEX");
    this.index = this.pinecone.Index(index);
  }

  async deleteMany(orgname: string, ids: string[]): Promise<void> {
    await this.index.namespace(orgname).deleteMany(ids);
  }

  async fetchAll(
    orgname: string,
    ids: string[]
  ): Promise<{
    vectors: {
      [vectorId: string]: number[];
    };
  }> {
    const vectors = {};
    const res = await this.index.namespace(orgname).fetch(ids);
    const records = res.records;
    for (const key in records) {
      if (records.hasOwnProperty(key)) {
        const record = records[key];
        vectors[record.id] = record.values;
      }
    }
    return { vectors };
  }

  async query(
    orgname: string,
    questionEmbedding: number[],
    topK: number,
    content?: { contentId: string }[]
  ): Promise<{ id: string; score: number }[]> {
    const queryResult = await this.index.namespace(orgname).query({
      includeMetadata: true,
      topK: topK,
      vector: questionEmbedding,
      ...(content.length > 0
        ? {
            filter: {
              contentId: {
                $in: content.map((doc) => doc.contentId),
              },
              orgname: orgname,
            },
          }
        : {
            filter: {
              orgname: orgname,
            },
          }),
    });
    return queryResult.matches.map((match) => ({
      id: match.id,
      score: match.score,
    }));
  }

  async upsert(
    orgname: string,
    contentId: string,
    embeddings: number[][]
  ): Promise<void> {
    for (let i = 0; i < embeddings.length; i += 100) {
      const endIndex = Math.min(i + 100, embeddings.length);
      const embeddingSlice = embeddings.slice(i, endIndex);
      await this.index.namespace(orgname).upsert(
        embeddingSlice.map((embedding, j) => {
          return {
            id: contentId + "__" + (i + j).toString(),
            metadata: {
              contentId: contentId,
              orgname: orgname,
            },
            values: embedding,
          };
        })
      );
    }
  }
}
