import { Injectable } from "@nestjs/common";
import { Member, Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateMemberDto } from "./dto/create-member.dto";
import { UpdateMemberDto } from "./dto/update-member.dto";

@Injectable()
export class MemberRepository extends BaseRepository<
  Member,
  CreateMemberDto,
  UpdateMemberDto,
  Prisma.MemberInclude,
  Prisma.MemberSelect
> {
  constructor(private prisma: PrismaService) {
    super(prisma.member);
  }

  async acceptMember(orgname: string, inviteEmail: string, username: string) {
    // Check org exists first just as added check
    await this.prisma.organization.findUniqueOrThrow({
      where: { orgname },
    });

    // Update Member
    return this.prisma.member.update({
      data: {
        inviteAccepted: true,
        user: {
          connect: {
            username,
          },
        },
      },
      where: {
        inviteEmail_orgname: {
          inviteEmail: inviteEmail,
          orgname: orgname,
        },
      },
    });
  }

  async create(orgname: string, createMemberDto: CreateMemberDto) {
    const existingUser = await this.prisma.user.findUnique({
      where: {
        email: createMemberDto.inviteEmail,
      },
    });
    return this.prisma.member.create({
      data: {
        organization: {
          connect: {
            orgname,
          },
        },
        ...(existingUser?.username
          ? { user: { connect: { username: existingUser.username } } }
          : {}),
        inviteEmail: createMemberDto.inviteEmail,
        role: createMemberDto.role,
      },
    });
  }

  async findByInviteEmail(orgname: string, inviteEmail: string) {
    return this.prisma.member.findUniqueOrThrow({
      where: {
        inviteEmail_orgname: {
          inviteEmail,
          orgname,
        },
      },
    });
  }
}
