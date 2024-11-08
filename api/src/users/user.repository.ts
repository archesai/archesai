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

@Injectable()
export class UserRepository extends BaseRepository<
  {
    authProviders: AuthProvider[];
    memberships: Member[];
  } & User,
  CreateUserDto,
  UpdateUserDto,
  Prisma.UserInclude,
  Prisma.UserSelect
> {
  constructor(private prisma: PrismaService) {
    super(prisma.user, {
      authProviders: true,
      memberships: true,
    });
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
      include: { authProviders: true, memberships: true },
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
      include: { authProviders: true, memberships: true },
    });
    return user;
  }

  async deactivate(id: string) {
    return this.prisma.user.update({
      data: {
        deactivated: true,
      },
      include: { authProviders: true, memberships: true },
      where: { id },
    });
  }

  async findOneByEmail(email: string) {
    return this.prisma.user.findUniqueOrThrow({
      include: { authProviders: true, memberships: true },
      where: { email },
    });
  }

  async setEmailVerified(id: string) {
    return this.prisma.user.update({
      data: {
        emailVerified: true,
      },
      include: { authProviders: true, memberships: true },
      where: { id },
    });
  }

  async setEmailVerifiedByEmail(email: string) {
    return this.prisma.user.update({
      data: {
        emailVerified: true,
      },
      include: { authProviders: true, memberships: true },
      where: { email },
    });
  }

  async updateEmail(id: string, email: string) {
    return this.prisma.user.update({
      data: {
        email,
      },
      include: { authProviders: true, memberships: true },
      where: { id },
    });
  }

  async updateRefreshToken(id: string, refreshToken: string) {
    return this.prisma.user.update({
      data: {
        refreshToken,
      },
      include: { authProviders: true, memberships: true },
      where: { id },
    });
  }
}
