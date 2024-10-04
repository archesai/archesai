import { Injectable } from "@nestjs/common";
import { AuthProviderType } from "@prisma/client";

import { PrismaService } from "../prisma/prisma.service";
import { UpdateUserDto } from "./dto/update-user.dto";
import { CreateUserInput } from "./types/create-user.type";

@Injectable()
export class UserRepository {
  constructor(private prisma: PrismaService) {}

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

  async create(createUser: CreateUserInput) {
    const prexistingMemberships = await this.prisma.member.findMany({
      where: {
        inviteEmail: createUser.email,
      },
    });
    const user = this.prisma.user.create({
      data: {
        ...createUser,
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

  async findAll() {
    return this.prisma.user.findMany({
      include: { authProviders: true, memberships: true },
      where: { deactivated: false },
    });
  }

  async findOne(id: string) {
    return this.prisma.user.findUniqueOrThrow({
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

  async remove(id: string) {
    await this.prisma.user.delete({ where: { id } });
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

  async update(id: string, updateUserDto: UpdateUserDto) {
    const user = await this.prisma.user.update({
      data: {
        firstName: updateUserDto.firstName,
        lastName: updateUserDto.lastName,
        username: updateUserDto.username,
      },
      include: { authProviders: true, memberships: true },
      where: { id },
    });
    return user;
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
