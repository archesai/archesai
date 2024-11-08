import { Injectable } from "@nestjs/common";
import { Content, Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateContentDto } from "./dto/create-content.dto";
import { UpdateContentDto } from "./dto/update-content.dto";

@Injectable()
export class ContentRepository extends BaseRepository<
  Content,
  CreateContentDto,
  UpdateContentDto,
  Prisma.ContentInclude,
  Prisma.ContentSelect
> {
  constructor(private prisma: PrismaService) {
    super(prisma.content);
  }

  async create(
    orgname: string,
    createContentDto: CreateContentDto,
    additionalData: {
      mimeType: string;
    }
  ) {
    return this.prisma.content.create({
      data: {
        ...createContentDto,
        mimeType: additionalData.mimeType,
        organization: {
          connect: {
            orgname,
          },
        },
      },
    });
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

  async incrementCredits(id: string, credits: number) {
    return this.prisma.content.update({
      data: {
        credits: {
          increment: credits,
        },
      },

      where: {
        id,
      },
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

  async removeMany(orgname: string, ids: string[]): Promise<void> {
    await this.prisma.content.deleteMany({
      where: {
        id: { in: ids },
        orgname,
      },
    });
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ContentUpdateInput) {
    return this.prisma.content.update({
      data: raw,
      where: {
        id,
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

    await this.prisma.content.createMany({
      data: records.map((record) => ({
        contentId,
        name: record.text.slice(0, 50),
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
