import { Injectable } from '@nestjs/common'
import { AuthProviderType, Prisma } from '@prisma/client'

import { BaseRepository } from '../common/base.repository'
import { PrismaService } from '../prisma/prisma.service'
import { CreateUserDto } from './dto/create-user.dto'

const USER_INCLUDE = {
  authProviders: true,
  memberships: true
} as const

@Injectable()
export class UserRepository extends BaseRepository<
  Prisma.UserDelegate,
  typeof USER_INCLUDE
> {
  constructor(private prisma: PrismaService) {
    super(prisma.user, USER_INCLUDE)
  }

  async addAuthProvider(
    email: string,
    provider: AuthProviderType,
    providerId: string
  ) {
    return await this.prisma.user.update({
      data: {
        authProviders: {
          create: {
            provider,
            providerId
          }
        }
      },
      include: USER_INCLUDE,
      where: { email }
    })
  }

  async create(createUserDto: CreateUserDto) {
    const prexistingMemberships = await this.prisma.member.findMany({
      where: {
        inviteEmail: createUserDto.email
      }
    })
    const user = this.prisma.user.create({
      data: {
        ...createUserDto,
        defaultOrgname: createUserDto.username,
        memberships: {
          connect: prexistingMemberships.map((m) => {
            return {
              inviteEmail_orgname: {
                inviteEmail: m.inviteEmail,
                orgname: m.orgname
              }
            }
          })
        }
      },
      include: USER_INCLUDE
    })
    return user
  }

  async deactivate(id: string) {
    await this.prisma.user.update({
      data: {
        deactivated: true
      },
      include: USER_INCLUDE,
      where: { id }
    })
  }

  async findOneByEmail(email: string) {
    return this.prisma.user.findUniqueOrThrow({
      include: USER_INCLUDE,
      where: { email }
    })
  }

  async findOneByUsername(username: string) {
    return this.prisma.user.findUniqueOrThrow({
      include: USER_INCLUDE,
      where: { username }
    })
  }
}
