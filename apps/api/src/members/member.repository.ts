import { Injectable } from '@nestjs/common'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { Prisma } from '@prisma/client'

@Injectable()
export class MemberRepository extends BaseRepository<Prisma.MemberDelegate> {
  constructor(private prisma: PrismaService) {
    super(prisma.member)
  }

  async findByInviteEmailAndOrgname(inviteEmail: string, orgname: string) {
    return this.prisma.member.findUniqueOrThrow({
      where: {
        inviteEmail_orgname: {
          inviteEmail,
          orgname
        }
      }
    })
  }

  async join(orgname: string, inviteEmail: string, username: string) {
    return this.prisma.member.update({
      data: {
        inviteAccepted: true,
        user: {
          connect: {
            username
          }
        }
      },
      where: {
        inviteEmail_orgname: {
          inviteEmail: inviteEmail,
          orgname: orgname
        }
      }
    })
  }
}
