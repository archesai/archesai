import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateContentDto } from "./dto/create-content.dto";
import { UpdateContentDto } from "./dto/update-content.dto";
import { ContentModel } from "./entities/content.entity";

@Injectable()
export class ContentRepository extends BaseRepository<
  ContentModel,
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
    const { labelIds, ...otherData } = createContentDto;
    return this.prisma.content.create({
      data: {
        ...otherData,
        labels: labelIds
          ? { connect: labelIds.map((id) => ({ id })) }
          : undefined,
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

  async update(
    orgname: string,
    contentId: string,
    updateContentDto: UpdateContentDto
  ) {
    const { labelIds, ...otherData } = updateContentDto;

    const data: Prisma.ContentUpdateInput = {
      ...otherData,
    };

    if (labelIds !== undefined) {
      data.labels = {
        set: labelIds.map((id) => ({ id })),
      };
    }
    return this.prisma.content.update({
      data: {
        ...otherData,
        labels: labelIds ? { set: labelIds.map((id) => ({ id })) } : undefined,
      },
      where: {
        id: contentId,
        orgname,
      },
    });
  }
}
