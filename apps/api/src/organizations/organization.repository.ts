import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'

@Injectable()
export class OrganizationRepository extends BaseRepository<Prisma.OrganizationDelegate> {
  constructor(private prisma: PrismaService) {
    super(prisma.organization)
  }

  async findByOrgname(orgname: string) {
    return this.prisma.organization.findFirstOrThrow({
      where: {
        orgname
      }
    })
  }

  async findByStripeCustomerId(stripeCustomerId: string) {
    return this.prisma.organization.findFirstOrThrow({
      where: {
        stripeCustomerId
      }
    })
  }
}
