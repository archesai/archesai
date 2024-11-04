import { Injectable, Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { CreateMemberDto } from "./dto/create-member.dto";
import { MemberQueryDto } from "./dto/member-query.dto";
import { UpdateMemberDto } from "./dto/update-member.dto";
import { MemberEntity } from "./entities/member.entity";
import { MemberRepository } from "./member.repository";

@Injectable()
export class MembersService
  implements
    BaseService<MemberEntity, CreateMemberDto, MemberQueryDto, UpdateMemberDto>
{
  private readonly logger = new Logger(MembersService.name);
  constructor(private memberRepository: MemberRepository) {}

  async acceptMember(orgname: string, inviteEmail: string, username: string) {
    this.logger.log(
      `Accepting member ${inviteEmail} to organization ${orgname}`
    );
    return new MemberEntity(
      await this.memberRepository.acceptMember(orgname, inviteEmail, username)
    );
  }

  async create(orgname: string, createMemberDto: CreateMemberDto) {
    const member = await this.memberRepository.create(orgname, createMemberDto);
    // await this.authService.sendEmailInvite(createMemberDto.inviteEmail); // FIXME
    return new MemberEntity(member);
  }

  async findAll(orgname: string, memberQueryDto: MemberQueryDto) {
    const { count, results } = await this.memberRepository.findAll(
      orgname,
      memberQueryDto
    );
    const memberEntities = results.map((member) => new MemberEntity(member));
    return new PaginatedDto<MemberEntity>({
      metadata: {
        limit: memberQueryDto.limit,
        offset: memberQueryDto.offset,
        totalResults: count,
      },
      results: memberEntities,
    });
  }

  async findById(id: string) {
    return new MemberEntity(await this.memberRepository.findById(id));
  }

  async findByInviteEmail(orgname: string, inviteEmail: string) {
    return new MemberEntity(
      await this.memberRepository.findByInviteEmail(orgname, inviteEmail)
    );
  }

  async remove(orgname: string, id: string) {
    await this.memberRepository.remove(orgname, id);
  }

  async update(orgname: string, id: string, updateMemberDto: UpdateMemberDto) {
    return new MemberEntity(
      await this.memberRepository.update(orgname, id, updateMemberDto)
    );
  }
}
