import { Injectable, Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { CreateMemberDto } from "./dto/create-member.dto";
import { UpdateMemberDto } from "./dto/update-member.dto";
import { MemberEntity, MemberModel } from "./entities/member.entity";
import { MemberRepository } from "./member.repository";

@Injectable()
export class MembersService extends BaseService<
  MemberEntity,
  CreateMemberDto,
  UpdateMemberDto,
  MemberRepository,
  MemberModel
> {
  private readonly logger = new Logger(MembersService.name);
  constructor(private memberRepository: MemberRepository) {
    super(memberRepository);
  }

  async join(orgname: string, inviteEmail: string, username: string) {
    this.logger.log(
      `Accepting member ${inviteEmail} to organization ${orgname}`
    );
    return this.toEntity(
      await this.memberRepository.join(orgname, inviteEmail, username)
    );
  }

  protected toEntity(model: MemberModel): MemberEntity {
    return new MemberEntity(model);
  }
}
