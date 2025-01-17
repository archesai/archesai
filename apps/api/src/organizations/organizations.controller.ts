import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post
} from '@nestjs/common'

import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { UserEntity } from '@/src/users/entities/user.entity'
import { CreateOrganizationDto } from '@/src/organizations/dto/create-organization.dto'
import { UpdateOrganizationDto } from '@/src/organizations/dto/update-organization.dto'
import { OrganizationEntity } from '@/src/organizations/entities/organization.entity'
import { OrganizationsService } from '@/src/organizations/organizations.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'
import { RoleTypeEnum } from '../members/entities/member.entity'

@Controller('/organizations')
export class OrganizationsController {
  constructor(private readonly organizationsService: OrganizationsService) {}

  /**
   * Create a new organization
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Authenticated([RoleTypeEnum.USER])
  @Post()
  async create(
    @Body() createOrganizationDto: CreateOrganizationDto,
    @CurrentUser() user: UserEntity
  ) {
    const organization = await this.organizationsService.create(
      createOrganizationDto
    )
    return this.organizationsService.addUserToOrganization(
      organization.orgname,
      user
    )
  }

  /**
   * Delete an organization
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Authenticated([RoleTypeEnum.ADMIN])
  @Delete(':orgname')
  async delete(@Param('orgname') orgname: string) {
    const organization = await this.organizationsService.findByOrgname(orgname)
    return this.organizationsService.remove(organization.id)
  }

  /**
   * Get an organization
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Authenticated([RoleTypeEnum.USER])
  @Get(':orgname')
  async findOne(@Param('orgname') orgname: string) {
    const organization = await this.organizationsService.findByOrgname(orgname)
    return new OrganizationEntity(
      await this.organizationsService.findOne(organization.id)
    )
  }

  /**
   * Update an organization
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Authenticated([RoleTypeEnum.ADMIN])
  @Patch(':orgname')
  async update(
    @Param('orgname') orgname: string,
    @Body() updateOrganizationDto: UpdateOrganizationDto
  ) {
    const organization = await this.organizationsService.findByOrgname(orgname)
    return new OrganizationEntity(
      await this.organizationsService.update(
        organization.id,
        updateOrganizationDto
      )
    )
  }
}
