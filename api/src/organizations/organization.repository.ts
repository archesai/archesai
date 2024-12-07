import { Injectable } from '@nestjs/common'
import { PlanType, Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { UserEntity } from '../users/entities/user.entity'
import { CreateOrganizationDto } from './dto/create-organization.dto'
import { UpdateOrganizationDto } from './dto/update-organization.dto'
import { OrganizationModel } from './entities/organization.entity'

@Injectable()
export class OrganizationRepository extends BaseRepository<
  OrganizationModel,
  CreateOrganizationDto,
  UpdateOrganizationDto,
  Prisma.OrganizationInclude,
  Prisma.OrganizationUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.organization)
  }

  async create(
    orgname: string,
    createOrganizationDto: CreateOrganizationDto,
    additionalData: {
      billingEnabled: boolean
      stripeCustomerId: string
      user: UserEntity
    }
  ) {
    const { billingEnabled, stripeCustomerId, user } = additionalData
    return this.prisma.organization.create({
      data: {
        ...createOrganizationDto,
        credits:
          // If this is their first org and their e-mail is verified, give them free credits
          // Otherwise, if billing is disabled, give them free credits
          billingEnabled ? (user.memberships?.length == 0 && user.emailVerified ? 0 : 0) : 100000000, // if this is their first org and their e-mail is verified, give them free credits
        // Add them as an admin to this organization
        members: {
          create: {
            inviteAccepted: true,
            inviteEmail: user.email, // FIXME
            role: 'ADMIN',
            user: {
              connect: {
                username: user.username
              }
            }
          }
        },
        plan: billingEnabled ? PlanType.FREE : PlanType.UNLIMITED,
        stripeCustomerId: stripeCustomerId
      }
    })
  }

  async findByOrgname(orgname: string) {
    return this.prisma.organization.findFirst({
      where: {
        orgname
      }
    })
  }

  async findByStripeCustomerId(stripeCustomerId: string) {
    return this.prisma.organization.findFirst({
      where: {
        stripeCustomerId
      }
    })
  }
}
