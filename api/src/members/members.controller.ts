import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  Query,
} from "@nestjs/common";
import {
  ApiBearerAuth,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from "@nestjs/swagger";

import {
  CurrentUser,
  CurrentUserDto,
} from "../auth/decorators/current-user.decorator";
import { BaseController } from "../common/base.controller";
import {
  ApiCrudOperation,
  Operation,
} from "../common/decorators/api-crud-operation.decorator";
import { SearchQueryDto } from "../common/dto/search-query.dto";
import { CreateMemberDto } from "./dto/create-member.dto";
import { UpdateMemberDto } from "./dto/update-member.dto";
import { MemberEntity } from "./entities/member.entity";
import { MembersService } from "./members.service";

@ApiBearerAuth()
@ApiTags("Organization - Members")
@Controller()
export class MembersController
  implements
    BaseController<
      MemberEntity,
      CreateMemberDto,
      SearchQueryDto,
      UpdateMemberDto
    >
{
  constructor(private readonly membersService: MembersService) {}

  @ApiCrudOperation(Operation.CREATE, "member", MemberEntity, true)
  @Post("/organizations/:orgname/members")
  async create(
    @Param("orgname") orgname: string,
    @Body() createMemberDto: CreateMemberDto
  ) {
    return this.membersService.create(orgname, createMemberDto);
  }

  @ApiCrudOperation(Operation.FIND_ALL, "member", MemberEntity, true)
  @Get("/organizations/:orgname/members")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() searchQueryDto: SearchQueryDto
  ) {
    return this.membersService.findAll(orgname, searchQueryDto);
  }

  @ApiOperation({
    description: "Accept invitation to this organization. ADMIN ONLY.",
    summary: "Accept invitation to this organization",
  })
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Email not verified", status: 403 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Member was successfully updated",
    status: 201,
    type: MemberEntity,
  })
  @Post("/organizations/:orgname/members/join")
  async join(
    @Param("orgname") orgname: string,
    @CurrentUser() user: CurrentUserDto
  ) {
    return this.membersService.join(orgname, user.email, user.username);
  }

  @ApiCrudOperation(Operation.DELETE, "member", MemberEntity, true)
  @Delete("/organizations/:orgname/members/:memberId")
  async remove(
    @Param("orgname") orgname: string,
    @Param("memberId") memberId: string
  ) {
    await this.membersService.remove(orgname, memberId);
  }

  @ApiCrudOperation(Operation.UPDATE, "member", MemberEntity, true)
  @Patch("/organizations/:orgname/members/:memberId")
  async update(
    @Param("orgname") orgname: string,
    @Param("memberId") memberId: string,
    @Body() updateMemberDto: UpdateMemberDto
  ) {
    return this.membersService.update(orgname, memberId, updateMemberDto);
  }
}
