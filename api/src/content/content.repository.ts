import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { CreateContentDto } from './dto/create-content.dto'
import { UpdateContentDto } from './dto/update-content.dto'
import { ContentModel } from './entities/content.entity'

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
}

@Injectable()
export class ContentRepository extends BaseRepository<
  ContentModel,
  CreateContentDto,
  UpdateContentDto,
  Prisma.ContentInclude,
  Prisma.ContentUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.content, CONTENT_INCLUDE)
  }

  async create(
    orgname: string,
    createContentDto: CreateContentDto,
    additionalData: {
      mimeType: string
    }
  ) {
    const { labels, ...otherData } = createContentDto
    return this.prisma.content.create({
      data: {
        ...otherData,
        labels:
          labels?.length > 0
            ? {
                connect: labels.map((name) => ({
                  name_orgname: {
                    name,
                    orgname
                  }
                }))
              }
            : undefined,
        mimeType: additionalData.mimeType,
        organization: {
          connect: {
            orgname
          }
        }
      },
      include: CONTENT_INCLUDE
    })
  }

  // Query vectors similar to a given embedding
  async query(
    orgname: string,
    embedding: number[],
    topK: number,
    contentIds?: string[]
  ): Promise<{ id: string; score: number }[]> {
    const results = await this.prisma.$queryRaw<{ id: string; score: number }[]>(Prisma.sql`
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

  async update(orgname: string, contentId: string, updateContentDto: UpdateContentDto) {
    const { labels, ...otherData } = updateContentDto

    return this.prisma.content.update({
      data: {
        ...otherData,
        labels:
          labels?.length > 0
            ? {
                set: labels.map((name) => ({
                  name_orgname: {
                    name,
                    orgname
                  }
                }))
              }
            : undefined
      },
      include: CONTENT_INCLUDE,
      where: {
        id: contentId,
        orgname
      }
    })
  }
}
