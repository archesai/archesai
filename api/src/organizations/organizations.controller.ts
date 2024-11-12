import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
} from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  CurrentUser,
  CurrentUserDto,
} from "../auth/decorators/current-user.decorator";
import { Roles } from "../auth/decorators/roles.decorator";
import {
  ApiCrudOperation,
  Operation,
} from "../common/decorators/api-crud-operation.decorator";
import { CreateOrganizationDto } from "./dto/create-organization.dto";
import { UpdateOrganizationDto } from "./dto/update-organization.dto";
import { OrganizationEntity } from "./entities/organization.entity";
import { OrganizationsService } from "./organizations.service";

@ApiBearerAuth()
@ApiTags("Organization")
@Roles("ADMIN")
@Controller("organizations")
export class OrganizationsController {
  constructor(private readonly organizationsService: OrganizationsService) {}

  @ApiCrudOperation(Operation.CREATE, "organization", OrganizationEntity, true)
  @Post()
  async create(
    @CurrentUser() user: CurrentUserDto,
    @Body() createOrganizationDto: CreateOrganizationDto
  ) {
    return new OrganizationEntity(
      await this.organizationsService.create(
        createOrganizationDto.orgname,
        createOrganizationDto,
        user
      )
    );
  }

  @ApiCrudOperation(Operation.DELETE, "organization", OrganizationEntity, true)
  @Delete(":orgname")
  async delete(@Param("orgname") orgname: string) {
    const organization = await this.organizationsService.findByOrgname(orgname);
    return this.organizationsService.remove(orgname, organization.id);
  }

  @ApiCrudOperation(Operation.GET, "organization", OrganizationEntity, true)
  @Get(":orgname")
  async findOne(@Param("orgname") orgname: string) {
    const organization = await this.organizationsService.findByOrgname(orgname);
    return new OrganizationEntity(
      await this.organizationsService.findOne(orgname, organization.id)
    );
  }

  @ApiCrudOperation(Operation.UPDATE, "organization", OrganizationEntity, true)
  @Patch(":orgname")
  async update(
    @Param("orgname") orgname: string,
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
