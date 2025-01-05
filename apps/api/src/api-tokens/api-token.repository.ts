import { Injectable } from '@nestjs/common'
import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { Prisma } from '@prisma/client'

@Injectable()
export class ApiTokenRepository extends BaseRepository<Prisma.ApiTokenDelegate> {
  constructor(private prisma: PrismaService) {
    super(prisma.apiToken)
  }
}
