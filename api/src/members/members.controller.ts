import { Controller, Param, Post, UseGuards } from '@nestjs/common'
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger'

import { CurrentUser } from '../auth/decorators/current-user.decorator'
import { BaseController } from '../common/base.controller'
import { UserEntity } from '../users/entities/user.entity'
import { CreateMemberDto } from './dto/create-member.dto'
import { UpdateMemberDto } from './dto/update-member.dto'
import { MemberEntity } from './entities/member.entity'
import { MembersService } from './members.service'
import { EmailVerifiedGuard } from '../auth/guards/email-verified.guard'

@ApiBearerAuth()
@ApiTags(`Members`)
@Controller('/organizations/:orgname/members')
export class MembersController extends BaseController<
  MemberEntity,
  CreateMemberDto,
  UpdateMemberDto,
  MembersService
>(MemberEntity, CreateMemberDto, UpdateMemberDto) {
  constructor(private readonly membersService: MembersService) {
    super(membersService)
  }

  /**
   * Accept invitation to this organization
   * @remarks This endpoint requires the user to be authenticated
   * @throws {401} Unauthorized
   * @throws {403} Email not verified
   * @throws {404} Not Found
   */
  @UseGuards(EmailVerifiedGuard)
  @Post('join')
  async join(
    @Param('orgname') orgname: string,
    @CurrentUser() currentUser: UserEntity
  ) {
    return this.membersService.join(
      orgname,
      currentUser.email,
      currentUser.username
    )
  }
}
