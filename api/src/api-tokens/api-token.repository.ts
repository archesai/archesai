import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { CreateApiTokenDto } from './dto/create-api-token.dto'
import { UpdateApiTokenDto } from './dto/update-api-token.dto'
import { ApiTokenModel } from './entities/api-token.entity'

@Injectable()
export class ApiTokenRepository extends BaseRepository<
  ApiTokenModel,
  CreateApiTokenDto,
  UpdateApiTokenDto,
  Prisma.ApiTokenInclude,
  Prisma.ApiTokenUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.apiToken)
  }

  async create(
    orgname: string,
    createApiTokenDto: CreateApiTokenDto,
    additionalData: {
      id: string
      key: string
      username: string
    }
  ) {
    return this.prisma.apiToken.create({
      data: {
        ...createApiTokenDto,
        id: additionalData.id,
        key: additionalData.key,
        organization: {
          connect: {
            orgname
          }
        },
        user: {
          connect: {
            username: additionalData.username
          }
        }
      }
    })
  }
}
