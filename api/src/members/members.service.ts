import { Injectable } from "@nestjs/common";
import { Member } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { CreateMemberDto } from "./dto/create-member.dto";
import { MemberQueryDto } from "./dto/member-query.dto";
import { UpdateMemberDto } from "./dto/update-member.dto";
import { MemberRepository } from "./member.repository";

@Injectable()
export class MembersService
  implements
    BaseService<Member, CreateMemberDto, MemberQueryDto, UpdateMemberDto>
{
  constructor(private memberRepository: MemberRepository) {}

  async acceptMember(orgname: string, inviteEmail: string, username: string) {
    return this.memberRepository.acceptMember(orgname, inviteEmail, username);
  }

  async create(orgname: string, createMemberDto: CreateMemberDto) {
    const member = await this.memberRepository.create(orgname, createMemberDto);
    // await this.authService.sendEmailInvite(createMemberDto.inviteEmail); // FIXME
    return member;
  }

  async findAll(orgname: string, memberQueryDto: MemberQueryDto) {
    return this.memberRepository.findAll(orgname, memberQueryDto);
  }

  async findById(id: string) {
    return this.memberRepository.findById(id);
  }

  async findByInviteEmail(orgname: string, inviteEmail: string) {
    return this.memberRepository.findByInviteEmail(orgname, inviteEmail);
  }

  async remove(orgname: string, id: string) {
    await this.memberRepository.remove(orgname, id);
  }

  async update(orgname: string, id: string, updateMemberDto: UpdateMemberDto) {
    return this.memberRepository.update(orgname, id, updateMemberDto);
  }
}
