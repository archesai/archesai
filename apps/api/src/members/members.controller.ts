import { Controller, Param, Post } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { BaseController } from '@/src/common/base.controller'
import { UserEntity } from '@/src/users/entities/user.entity'
import { CreateMemberDto } from '@/src/members/dto/create-member.dto'
import { UpdateMemberDto } from '@/src/members/dto/update-member.dto'
import { MemberEntity } from '@/src/members/entities/member.entity'
import { MembersService } from '@/src/members/members.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@ApiTags(`Members`)
@Authenticated()
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
