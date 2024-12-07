import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { CreateLabelDto } from './dto/create-label.dto'
import { UpdateLabelDto } from './dto/update-label.dto'
import { LabelModel } from './entities/label.entity'

@Injectable()
export class LabelRepository extends BaseRepository<
  LabelModel,
  CreateLabelDto,
  UpdateLabelDto,
  Prisma.LabelInclude,
  Prisma.LabelUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.label)
  }
}
