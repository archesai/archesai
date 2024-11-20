import { Body, Controller, Param, Patch, Post } from "@nestjs/common";
import { ApiBearerAuth } from "@nestjs/swagger";

import { CurrentUser } from "../auth/decorators/current-user.decorator";
import { BaseController } from "../common/base.controller";
import { UserEntity } from "../users/entities/user.entity";
import { CreateOrganizationDto } from "./dto/create-organization.dto";
import { UpdateOrganizationDto } from "./dto/update-organization.dto";
import { OrganizationEntity } from "./entities/organization.entity";
import { OrganizationsService } from "./organizations.service";

@ApiBearerAuth()
@Controller("/organizations")
export class OrganizationsController extends BaseController<
  OrganizationEntity,
  CreateOrganizationDto,
  UpdateOrganizationDto,
  OrganizationsService
> {
  constructor(private readonly organizationsService: OrganizationsService) {
    super(organizationsService);
  }

  /**
   * Create a new organization
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Post()
  async create(
    @Param("orgname") orgname: string,
    @Body() createOrganizationDto: CreateOrganizationDto,
    @CurrentUser() user: UserEntity
  ) {
    return new OrganizationEntity(
      await this.organizationsService.create(
        createOrganizationDto.orgname,
        createOrganizationDto,
        user
      )
    );
  }

  /**
   * Update an organization
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Patch(":orgname")
  async update(
    @Param("orgname") orgname: string,
    @Param("id") id: string,
    @Body() updateOrganizationDto: UpdateOrganizationDto
  ) {
    const organization = await this.organizationsService.findByOrgname(orgname);
    return new OrganizationEntity(
      await this.organizationsService.update(
        orgname,
        organization.id,
        updateOrganizationDto
      )
    );
  }
}
