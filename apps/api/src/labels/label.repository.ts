import { Injectable } from '@nestjs/common'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { Prisma } from '@prisma/client'

@Injectable()
export class LabelRepository extends BaseRepository<Prisma.LabelDelegate> {
  constructor(private prisma: PrismaService) {
    super(prisma.label, {})
  }
}
