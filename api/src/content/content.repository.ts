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
  Prisma.ContentSelect,
  Prisma.ContentUpdateInput
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
}
