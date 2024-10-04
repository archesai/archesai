import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { PrismaService } from "../prisma/prisma.service";
import { BaseVectorDBService, VectorDBService } from "./vector-db.service";

@Injectable()
export class PgVectorDBService
  extends BaseVectorDBService
  implements VectorDBService
{
  constructor(private prisma: PrismaService) {
    super();
  }

  // Delete specific vectors by their IDs
  async deleteMany(orgname: string, ids: string[]): Promise<void> {
    await this.prisma.vectorRecord.deleteMany({
      where: {
        id: { in: ids },
        orgname,
      },
    });
  }

  // Fetch vectors by their IDs
  async fetchAll(
    orgname: string,
    ids: string[]
  ): Promise<{
    vectors: {
      [vectorId: string]: number[];
    };
  }> {
    const vectors = {};

    // Use raw SQL query because 'embedding' is an unsupported type
    const records = await this.prisma.$queryRaw<
      Array<{ embedding: number[]; id: string }>
    >(Prisma.sql`
      SELECT id, embedding::float8[] AS embedding
      FROM "VectorRecord"
      WHERE orgname = ${orgname} AND id IN (${Prisma.join(ids)})
    `);

    for (const record of records) {
      vectors[record.id] = record.embedding;
    }

    return { vectors };
  }

  // Query vectors similar to a given embedding
  async query(
    orgname: string,
    questionEmbedding: number[],
    topK: number,
    contents?: { contentId: string }[]
  ): Promise<{ id: string; score: number }[]> {
    const contentIds = contents?.map((content) => content.contentId);

    const results = await this.prisma.$queryRaw<
      { id: string; score: number }[]
    >(Prisma.sql`
      SELECT
        id,
        1 - (embedding <#> ${questionEmbedding}::vector) AS score
      FROM
        "VectorRecord"
      WHERE
        orgname = ${orgname}
        ${
          contentIds?.length
            ? Prisma.sql`AND "contentId" IN (${Prisma.join(contentIds)})`
            : Prisma.empty
        }
      ORDER BY
        embedding <#> ${questionEmbedding}::vector ASC
      LIMIT ${topK};
    `);

    return results;
  }

  // Upsert vectors with embeddings and text
  async upsert(
    orgname: string,
    contentId: string,
    embeddings: number[][],
    texts?: string[] // Ensure this array matches the embeddings array length
  ): Promise<void> {
    const values = embeddings.map((embedding, index) => ({
      contentId,
      embedding,
      id: `${contentId}__${index}`,
      orgname,
      text: texts ? texts[index] : "", // Include text data
    }));

    // Construct the VALUES clause using Prisma.sql and Prisma.join
    const valuesSql = Prisma.join(
      values.map(
        (record) =>
          Prisma.sql`(
            ${record.id},
            ${record.orgname},
            ARRAY[${Prisma.join(record.embedding)}]::vector,
            ${record.contentId},
            ${record.text},
            NOW()
          )`
      ),
      `, `
    );

    // Construct the full query using Prisma.sql tagged template
    const query = Prisma.sql`
      INSERT INTO "VectorRecord" (id, orgname, embedding, "contentId", text, "updatedAt")
      VALUES ${valuesSql}
      ON CONFLICT (id) DO UPDATE SET
        orgname = EXCLUDED.orgname,
        embedding = EXCLUDED.embedding,
        "contentId" = EXCLUDED."contentId",
        text = EXCLUDED.text,
        "updatedAt" = NOW();
    `;

    // Execute the query using $executeRaw with the Prisma.sql object
    await this.prisma.$executeRaw(query);
  }
}
