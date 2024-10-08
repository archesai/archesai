import { Injectable } from "@nestjs/common";
import { VectorRecord } from "@prisma/client";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { VectorRecordQueryDto } from "./dto/vector-record-query.dto";

@Injectable()
export class VectorRecordRepository
  implements
    BaseRepository<VectorRecord, undefined, VectorRecordQueryDto, undefined>
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
      FROM "VectorRecord"
      WHERE orgname = ${orgname} AND id IN (${Prisma.join(ids)})
    `);

    for (const record of records) {
      vectors[record.id] = record.embedding;
    }

    return { vectors };
  }

  async findAll(
    orgname: string,
    vectorRecordQueryDto: VectorRecordQueryDto,
    contentId?: string
  ) {
    const count = await this.prisma.vectorRecord.count({
      where: {
        contentId,
        orgname,
        ...(vectorRecordQueryDto.searchTerm
          ? {
              OR: [{ text: { contains: vectorRecordQueryDto.searchTerm } }],
            }
          : undefined),
      },
    });
    const vectorRecords = await this.prisma.vectorRecord.findMany({
      orderBy: {
        [vectorRecordQueryDto.sortBy || "createdAt"]:
          vectorRecordQueryDto.sortDirection,
      },
      skip: vectorRecordQueryDto.offset,
      take: vectorRecordQueryDto.limit,
      where: {
        contentId,
        orgname,
        ...(vectorRecordQueryDto.searchTerm
          ? {
              OR: [{ text: { contains: vectorRecordQueryDto.searchTerm } }],
            }
          : undefined),
      },
    });
    return { count, results: vectorRecords };
  }

  async findOne(id: string) {
    return this.prisma.vectorRecord.findUniqueOrThrow({
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
        "VectorRecord"
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
    await this.prisma.vectorRecord.delete({
      where: {
        id,
        orgname,
      },
    });
  }

  // Remove vectors by their IDs
  async removeMany(orgname: string, ids: string[]): Promise<void> {
    await this.prisma.vectorRecord.deleteMany({
      where: {
        id: { in: ids },
        orgname,
      },
    });
  }

  // Upsert vectors with embeddings and text
  async upsert(
    orgname: string,
    contentId: string,
    records: {
      embedding: number[];
      text: string;
    }[]
  ): Promise<void> {
    // Construct the VALUES clause using Prisma.sql and Prisma.join
    const valuesSql = Prisma.join(
      records.map(
        (record, i) =>
          Prisma.sql`(
            ${contentId + "_" + i},
            ${orgname},
            ARRAY[${Prisma.join(record.embedding)}]::vector,
            ${contentId},
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
