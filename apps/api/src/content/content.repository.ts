import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'

const CONTENT_INCLUDE = {
  children: {
    select: {
      id: true,
      name: true
    }
  },
  consumedBy: {
    select: {
      id: true,
      name: true
    }
  },
  labels: true,
  parent: {
    select: {
      id: true,
      name: true
    }
  },
  producedBy: {
    select: {
      id: true,
      name: true
    }
  }
} as const

@Injectable()
export class ContentRepository extends BaseRepository<
  Prisma.ContentDelegate,
  typeof CONTENT_INCLUDE
> {
  constructor(private prisma: PrismaService) {
    super(prisma.content, CONTENT_INCLUDE)
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
        ${contentIds?.length ? Prisma.sql`AND "contentId" IN (${Prisma.join(contentIds)})` : Prisma.empty}
      ORDER BY
        embedding <#> ${embedding}::vector ASC
      LIMIT ${topK};
    `)

    return results
  }
}
