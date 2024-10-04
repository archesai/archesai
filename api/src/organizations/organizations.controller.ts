import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  // Delete,
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
import { Roles } from "../auth/decorators/roles.decorator";
import { CreateOrganizationDto } from "./dto/create-organization.dto";
import { UpdateOrganizationDto } from "./dto/update-organization.dto";
import { OrganizationEntity } from "./entities/organization.entity";
import { OrganizationsService } from "./organizations.service";

@Roles("ADMIN")
@ApiBearerAuth()
@ApiTags("Organization")
@Controller("organizations")
export class OrganizationsController {
  constructor(private readonly organizationsService: OrganizationsService) {}

  @Post()
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({
    description: "Organization was successfully created",
    status: 201,
    type: OrganizationEntity,
  })
  @ApiResponse({
    description: "Email not verified",
    status: 403,
  })
  @ApiOperation({
    description: "Create an organization. ADMIN ONLY.",
    summary: "Create an organization",
  })
  async create(
    @CurrentUser() user: CurrentUserDto,
    @Body() createOrganizationDto: CreateOrganizationDto
  ) {
    return new OrganizationEntity(
      await this.organizationsService.createAndInitialize(
        user,
        createOrganizationDto
      )
    );
  }

  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Organization was successfully deleted",
    status: 200,
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @Delete(":orgname")
  @ApiOperation({
    description: "Delete an organization. ADMIN ONLY.",
    summary: "Delete an organization",
  })
  async delete(@Param("orgname") organization: string) {
    return this.organizationsService.removeByName(organization);
  }

  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Organization was successfully returned",
    status: 200,
    type: OrganizationEntity,
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @Get(":orgname")
  @ApiOperation({
    description: "Get an organization. ADMIN ONLY.",
    summary: "Get an organization",
  })
  async findOne(@Param("orgname") organization: string) {
    return new OrganizationEntity(
      await this.organizationsService.findOneByName(organization)
    );
  }

  @Patch(":orgname")
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Organization was successfully updated",
    status: 200,
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @ApiOperation({
    description: "Update an organization. ADMIN ONLY.",
    summary: "Update an organization",
  })
  async update(
    @Param("orgname") organization: string,
    @Body() updateOrganizationDto: UpdateOrganizationDto
  ) {
    return new OrganizationEntity(
      await this.organizationsService.updateByName(
        organization,
        updateOrganizationDto
      )
    );
  }
}
