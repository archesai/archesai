import { Injectable } from '@nestjs/common'
import { Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { CreateMemberDto } from './dto/create-member.dto'
import { UpdateMemberDto } from './dto/update-member.dto'
import { MemberModel } from './entities/member.entity'

@Injectable()
export class MemberRepository extends BaseRepository<
  MemberModel,
  CreateMemberDto,
  UpdateMemberDto,
  Prisma.MemberInclude,
  Prisma.MemberUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.member)
  }

  async create(orgname: string, createMemberDto: CreateMemberDto) {
    const existingUser = await this.prisma.user.findUnique({
      where: {
        email: createMemberDto.inviteEmail
      }
    })
    return this.prisma.member.create({
      data: {
        organization: {
          connect: {
            orgname
          }
        },
        ...(existingUser?.username
          ? { user: { connect: { username: existingUser.username } } }
          : {}),
        inviteEmail: createMemberDto.inviteEmail,
        role: createMemberDto.role
      }
    })
  }

  async join(orgname: string, inviteEmail: string, username: string) {
    await this.prisma.organization.findUniqueOrThrow({
      where: { orgname }
    })

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
