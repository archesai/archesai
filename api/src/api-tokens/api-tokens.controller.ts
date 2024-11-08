import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
} from "@nestjs/common";
import { Query } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  CurrentUser,
  CurrentUserDto,
} from "../auth/decorators/current-user.decorator";
import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { SearchQueryDto } from "../common/search-query";
import { ApiTokensService } from "./api-tokens.service";
import { CreateApiTokenDto } from "./dto/create-api-token.dto";
import { UpdateApiTokenDto } from "./dto/update-api-token.dto";
import { ApiTokenEntity } from "./entities/api-token.entity";

@ApiBearerAuth()
@ApiTags("API Tokens")
@Controller()
export class ApiTokensController
  implements
    BaseController<
      ApiTokenEntity,
      CreateApiTokenDto,
      SearchQueryDto,
      UpdateApiTokenDto
    >
{
  constructor(private readonly apiTokensService: ApiTokensService) {}

  @ApiCrudOperation(Operation.CREATE, "API token", ApiTokenEntity, true)
  @Post("/organizations/:orgname/api-tokens")
  async create(
    @Param("orgname") orgname: string,
    @Body() createTokenDto: CreateApiTokenDto,
    @CurrentUser() currentUserDto?: CurrentUserDto
  ) {
    return this.apiTokensService.create(
      orgname,
      createTokenDto,
      currentUserDto.id
    );
  }

  @ApiCrudOperation(Operation.FIND_ALL, "API token", ApiTokenEntity, true)
  @Get("/organizations/:orgname/api-tokens")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() searchQueryDto: SearchQueryDto
  ) {
    return this.apiTokensService.findAll(orgname, searchQueryDto);
  }

  @ApiCrudOperation(Operation.GET, "API token", ApiTokenEntity, true)
  @Get("/organizations/:orgname/api-tokens/:id")
  async findOne(@Param("orgname") orgname: string, @Param("id") id: string) {
    return this.apiTokensService.findOne(orgname, id);
  }

  @ApiCrudOperation(Operation.DELETE, "API token", ApiTokenEntity, true)
  @Delete("/organizations/:orgname/api-tokens/:id")
  async remove(@Param("orgname") orgname: string, @Param("id") id: string) {
    await this.apiTokensService.remove(orgname, id);
  }

  @ApiCrudOperation(Operation.UPDATE, "API token", ApiTokenEntity, true)
  @Patch("/organizations/:orgname/api-tokens/:id")
  async update(
    @Param("orgname") orgname: string,
    @Param("id") id: string,
    @Body() updateApiTokenDto: UpdateApiTokenDto
  ) {
    return this.apiTokensService.update(orgname, id, updateApiTokenDto);
  }
}
