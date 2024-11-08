import { Injectable } from "@nestjs/common";
import {
  AuthProvider,
  AuthProviderType,
  Member,
  Prisma,
  User,
} from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateUserDto } from "./dto/create-user.dto";
import { UpdateUserDto } from "./dto/update-user.dto";

const USER_INCLUDE = {
  authProviders: true,
  memberships: true,
};

@Injectable()
export class UserRepository extends BaseRepository<
  {
    authProviders: AuthProvider[];
    memberships: Member[];
  } & User,
  CreateUserDto,
  UpdateUserDto,
  Prisma.UserInclude,
  Prisma.UserSelect,
  Prisma.UserUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.user, USER_INCLUDE);
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
            providerId,
          },
        },
      },
      include: USER_INCLUDE,
      where: { email },
    });
  }

  async create(orgname: string, createUserDto: CreateUserDto) {
    const prexistingMemberships = await this.prisma.member.findMany({
      where: {
        inviteEmail: createUserDto.email,
      },
    });
    const user = this.prisma.user.create({
      data: {
        ...createUserDto,
        defaultOrgname: createUserDto.username,
        memberships: {
          connect: prexistingMemberships.map((m) => {
            return {
              inviteEmail_orgname: {
                inviteEmail: m.inviteEmail,
                orgname: m.orgname,
              },
            };
          }),
        },
      },
      include: USER_INCLUDE,
    });
    return user;
  }

  async deactivate(id: string) {
    return this.prisma.user.update({
      data: {
        deactivated: true,
      },
      include: USER_INCLUDE,
      where: { id },
    });
  }

  async findOneByEmail(email: string) {
    return this.prisma.user.findUniqueOrThrow({
      include: USER_INCLUDE,
      where: { email },
    });
  }
}
