import { Injectable } from "@nestjs/common";
import { TextChunk } from "@prisma/client";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { TextChunkQueryDto } from "./dto/text-chunk-query.dto";

@Injectable()
export class TextChunkRepository
  implements BaseRepository<TextChunk, undefined, TextChunkQueryDto, undefined>
{
  constructor(private prisma: PrismaService) {}

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
      FROM "TextChunk"
      WHERE 
        orgname = ${orgname}
        ${
          ids.length
            ? Prisma.sql`AND id IN (${Prisma.join(ids)})`
            : Prisma.empty
        }
    `);

    for (const record of records) {
      vectors[record.id] = record.embedding;
    }

    return { vectors };
  }

  async findAll(
    orgname: string,
    textChunkRepository: TextChunkQueryDto,
    contentId?: string
  ) {
    const count = await this.prisma.textChunk.count({
      where: {
        contentId,
        orgname,
        ...(textChunkRepository.searchTerm
          ? {
              OR: [{ text: { contains: textChunkRepository.searchTerm } }],
            }
          : undefined),
      },
    });
    const textChunks = await this.prisma.textChunk.findMany({
      orderBy: {
        [textChunkRepository.sortBy || "createdAt"]:
          textChunkRepository.sortDirection,
      },
      skip: textChunkRepository.offset,
      take: textChunkRepository.limit,
      where: {
        contentId,
        orgname,
        ...(textChunkRepository.searchTerm
          ? {
              OR: [{ text: { contains: textChunkRepository.searchTerm } }],
            }
          : undefined),
      },
    });
    return { count, results: textChunks };
  }

  async findOne(id: string) {
    return this.prisma.textChunk.findUniqueOrThrow({
      where: { id },
    });
  }

  // Query vectors similar to a given embedding
  async query(
    orgname: string,
    embedding: number[],
    topK: number,
    contentIds?: string[]
  ): Promise<{ id: string; score: number }[]> {
    const results = await this.prisma.$queryRaw<
      { id: string; score: number }[]
    >(Prisma.sql`
      SELECT
        id,
        1 - (embedding <#> ${embedding}::vector) AS score
      FROM
        "TextChunk"
      WHERE
        orgname = ${orgname}
        ${
          contentIds?.length
            ? Prisma.sql`AND "contentId" IN (${Prisma.join(contentIds)})`
            : Prisma.empty
        }
      ORDER BY
        embedding <#> ${embedding}::vector ASC
      LIMIT ${topK};
    `);

    return results;
  }

  // Remove vector by id
  async remove(orgname: string, id: string): Promise<void> {
    await this.prisma.textChunk.delete({
      where: {
        id,
        orgname,
      },
    });
  }

  // Remove vectors by their IDs
  async removeMany(orgname: string, ids: string[]): Promise<void> {
    await this.prisma.textChunk.deleteMany({
      where: {
        id: { in: ids },
        orgname,
      },
    });
  }

  async upsertTextChunks(
    orgname: string,
    contentId: string,
    records: {
      text: string;
    }[]
  ): Promise<void> {
    if (!records.length) {
      return;
    }

    await this.prisma.textChunk.createMany({
      data: records.map((record) => ({
        contentId,
        orgname,
        text: record.text,
      })),
    });
  }

  // Upsert vectors with embeddings and text
  async upsertVectors(
    orgname: string,
    contentId: string,
    records: {
      embedding: number[];
      textChunkId: string;
    }[]
  ): Promise<void> {
    // Construct the VALUES clause using Prisma.sql and Prisma.join
    if (!records.length) {
      return;
    }
    const valuesSql = Prisma.join(
      records.map(
        (record) =>
          Prisma.sql`(
            ${record.textChunkId},
            ${orgname},
            ARRAY[${Prisma.join(record.embedding)}]::vector,
            ${contentId},
            NOW()
          )`
      ),
      `, `
    );

    // Construct the full query using Prisma.sql tagged template
    const query = Prisma.sql`
      INSERT INTO "TextChunk" (id, orgname, embedding, "contentId", "updatedAt")
      VALUES ${valuesSql}
      ON CONFLICT (id) DO UPDATE SET
        orgname = EXCLUDED.orgname,
        embedding = EXCLUDED.embedding,
        "contentId" = EXCLUDED."contentId",
        "updatedAt" = NOW();
    `;

    // Execute the query using $executeRaw with the Prisma.sql object
    await this.prisma.$executeRaw(query);
  }
}
