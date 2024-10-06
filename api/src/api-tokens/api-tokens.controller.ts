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

import { Roles } from "../auth/decorators/roles.decorator";
import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { PaginatedDto } from "../common/paginated.dto";
import { ApiTokensService } from "./api-tokens.service";
import { ApiTokenQueryDto } from "./dto/api-token-query.dto";
import { CreateApiTokenDto } from "./dto/create-api-token.dto";
import { UpdateApiTokenDto } from "./dto/update-api-token.dto";
import { ApiTokenEntity } from "./entities/api-token.entity";

@ApiBearerAuth()
@ApiTags("Organization - API Tokens")
@Roles("ADMIN")
@Controller()
export class ApiTokensController
  implements
    BaseController<
      ApiTokenEntity,
      CreateApiTokenDto,
      ApiTokenQueryDto,
      UpdateApiTokenDto
    >
{
  constructor(private readonly apiTokensService: ApiTokensService) {}

  @ApiCrudOperation(Operation.CREATE, "API token", ApiTokenEntity, true)
  @Post("/organizations/:orgname/api-tokens")
  async create(
    @Param("orgname") orgname: string,
    @Body() createTokenDto: CreateApiTokenDto
  ) {
    const apiToken = await this.apiTokensService.create(
      orgname,
      createTokenDto
    );
    return new ApiTokenEntity(apiToken);
  }

  @ApiCrudOperation(Operation.FIND_ALL, "API token", ApiTokenEntity, true)
  @Get("/organizations/:orgname/api-tokens")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() apiTokenQueryDto: ApiTokenQueryDto
  ) {
    const { count, results } = await this.apiTokensService.findAll(
      orgname,
      apiTokenQueryDto
    );

    return new PaginatedDto<ApiTokenEntity>({
      metadata: {
        limit: 10,
        offset: 0,
        totalResults: count,
      },
      results: results.map((val) => new ApiTokenEntity(val)),
    });
  }

  @ApiCrudOperation(Operation.GET, "API token", ApiTokenEntity, true)
  @Get("/organizations/:orgname/api-tokens/:id")
  async findOne(@Param("orgname") orgname: string, @Param("id") id: string) {
    const apiToken = await this.apiTokensService.findOne(orgname, id);
    return new ApiTokenEntity(apiToken);
  }

  @ApiCrudOperation(Operation.DELETE, "API token", ApiTokenEntity, true)
  @Delete("/organizations/:orgname/api-tokens/:id")
  async remove(@Param("orgname") orgname: string, @Param("id") id: string) {
    return this.apiTokensService.remove(orgname, id);
  }

  @ApiCrudOperation(Operation.UPDATE, "API token", ApiTokenEntity, true)
  @Patch("/organizations/:orgname/api-tokens/:id")
  async update(
    @Param("orgname") orgname: string,
    @Param("id") id: string,
    @Body() updateApiTokenDto: UpdateApiTokenDto
  ) {
    const apiToken = await this.apiTokensService.update(
      orgname,
      id,
      updateApiTokenDto
    );
    return new ApiTokenEntity(apiToken);
  }
}
