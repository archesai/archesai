import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Query,
} from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { PaginatedDto } from "../common/paginated.dto";
import { ToolQueryDto } from "./dto/tool-query.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolEntity } from "./entities/tool.entity";
import { ToolsService } from "./tools.service";

@ApiBearerAuth()
@ApiTags("Tools")
@Controller("organizations/:orgname/tools")
export class ToolsController
  implements BaseController<ToolEntity, undefined, ToolQueryDto, UpdateToolDto>
{
  constructor(private readonly toolsService: ToolsService) {}

  @ApiCrudOperation(Operation.FIND_ALL, "tool", ToolEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() toolQueryDto: ToolQueryDto
  ) {
    const { count, results } = await this.toolsService.findAll(
      orgname,
      toolQueryDto
    );
    const toolsEntities = await Promise.all(
      results.map(async (tools) => {
        const populated = await this.toolsService.populateReadUrl(tools);
        return new ToolEntity(populated);
      })
    );
    return new PaginatedDto<ToolEntity>({
      metadata: {
        limit: toolQueryDto.limit,
        offset: toolQueryDto.offset,
        totalResults: count,
      },
      results: toolsEntities,
    });
  }

  @ApiCrudOperation(Operation.GET, "tool", ToolEntity, true)
  @Get("/:toolId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string
  ) {
    return new ToolEntity(await this.toolsService.findOne(toolId));
  }

  @ApiCrudOperation(Operation.DELETE, "tools", ToolEntity, true)
  @Delete("/:toolId")
  remove(@Param("orgname") orgname: string, @Param("toolId") toolId: string) {
    return this.toolsService.remove(orgname, toolId);
  }

  @ApiCrudOperation(Operation.UPDATE, "tools", ToolEntity, true)
  @Patch("/:toolId")
  async update(
    @Param("orgname") orgname: string,
    @Param("toolId") toolId: string,
    @Body() updateContentDto: UpdateToolDto
  ) {
    return new ToolEntity(
      await this.toolsService.update(orgname, toolId, updateContentDto)
    );
  }
}
