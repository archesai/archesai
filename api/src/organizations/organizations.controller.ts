import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post
} from '@nestjs/common'
import { ApiBearerAuth } from '@nestjs/swagger'

import { CurrentUser } from '../auth/decorators/current-user.decorator'
import { UserEntity } from '../users/entities/user.entity'
import { CreateOrganizationDto } from './dto/create-organization.dto'
import { UpdateOrganizationDto } from './dto/update-organization.dto'
import { OrganizationEntity } from './entities/organization.entity'
import { OrganizationsService } from './organizations.service'

@ApiBearerAuth()
@Controller('/organizations')
export class OrganizationsController {
  constructor(private readonly organizationsService: OrganizationsService) {}

  /**
   * Create a new organization
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Post()
  async create(
    @Body() createOrganizationDto: CreateOrganizationDto,
    @CurrentUser() user: UserEntity
  ) {
    return new OrganizationEntity(
      await this.organizationsService.create(
        createOrganizationDto.orgname,
        createOrganizationDto,
        user
      )
    )
  }

  /**
   * Delete an organization
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Delete(':orgname')
  async delete(@Param('orgname') orgname: string) {
    const organization = await this.organizationsService.findByOrgname(orgname)
    return this.organizationsService.remove(orgname, organization.id)
  }

  /**
   * Get an organization
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Get(':orgname')
  async findOne(@Param('orgname') orgname: string) {
    const organization = await this.organizationsService.findByOrgname(orgname)
    return new OrganizationEntity(
      await this.organizationsService.findOne(orgname, organization.id)
    )
  }

  /**
   * Update an organization
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Patch(':orgname')
  async update(
    @Param('orgname') orgname: string,
    @Body() updateOrganizationDto: UpdateOrganizationDto
  ) {
    const organization = await this.organizationsService.findByOrgname(orgname)
    return new OrganizationEntity(
      await this.organizationsService.update(
        orgname,
        organization.id,
        updateOrganizationDto
      )
    )
  }
}
