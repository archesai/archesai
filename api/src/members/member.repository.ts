import { Injectable } from "@nestjs/common";
import { Member } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateMemberDto } from "./dto/create-member.dto";
import { MemberQueryDto } from "./dto/member-query.dto";
import { UpdateMemberDto } from "./dto/update-member.dto";

@Injectable()
export class MemberRepository
  implements
    BaseRepository<Member, CreateMemberDto, MemberQueryDto, UpdateMemberDto>
{
  constructor(private prisma: PrismaService) {}

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

  async findAll(orgname: string, memberQueryDto: MemberQueryDto) {
    const whereConditions = {
      createdAt: {
        gte: memberQueryDto.startDate,
        lte: memberQueryDto.endDate,
      },
      orgname,
    };
    if (memberQueryDto.filters) {
      memberQueryDto.filters.forEach((filter) => {
        whereConditions[filter.field] = { contains: filter.value };
      });
    }
    const count = await this.prisma.member.count({ where: whereConditions });
    const members = await this.prisma.member.findMany({
      orderBy: {
        [memberQueryDto.sortBy]: memberQueryDto.sortDirection,
      },
      skip: memberQueryDto.offset,
      take: memberQueryDto.limit,
      where: whereConditions,
    });
    return { count, results: members };
  }

  async findById(id: string) {
    return this.prisma.member.findUniqueOrThrow({
      where: {
        id,
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

  async remove(orgname: string, id: string) {
    await this.prisma.member.delete({ where: { id } });
  }
  async update(orgname: string, id: string, updateMemberDto: UpdateMemberDto) {
    return this.prisma.member.update({
      data: {
        inviteEmail: updateMemberDto.inviteEmail,
        role: updateMemberDto.role,
      },
      where: {
        id,
      },
    });
  }
}
