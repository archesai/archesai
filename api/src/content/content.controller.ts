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
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { PaginatedDto } from "../common/paginated.dto";
import { ContentService } from "./content.service";
import { ContentQueryDto } from "./dto/content-query.dto";
import { CreateContentDto } from "./dto/create-content.dto";
import { UpdateContentDto } from "./dto/update-content.dto";
import { ContentEntity } from "./entities/content.entity";

@ApiBearerAuth()
@ApiTags("Content")
@Controller("organizations/:orgname/content")
export class ContentController
  implements
    BaseController<ContentEntity, undefined, ContentQueryDto, UpdateContentDto>
{
  constructor(private readonly contentService: ContentService) {}

  @ApiCrudOperation(Operation.CREATE, "content", ContentEntity, true)
  @Post()
  async create(
    @Param("orgname") orgname: string,
    @Body() createContentDto: CreateContentDto
  ) {
    return new ContentEntity(
      await this.contentService.create(orgname, createContentDto)
    );
  }

  @ApiCrudOperation(Operation.FIND_ALL, "content", ContentEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() contentQueryDto: ContentQueryDto
  ) {
    const { count, results } = await this.contentService.findAll(
      orgname,
      contentQueryDto
    );
    const contentEntities = await Promise.all(
      results.map(async (content) => {
        const populated = await this.contentService.populateReadUrl(content);
        return new ContentEntity(populated);
      })
    );
    return new PaginatedDto<ContentEntity>({
      metadata: {
        limit: contentQueryDto.limit,
        offset: contentQueryDto.offset,
        totalResults: count,
      },
      results: contentEntities,
    });
  }

  @ApiCrudOperation(Operation.GET, "content", ContentEntity, true)
  @Get("/:contentId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("contentId") contentId: string
  ) {
    return new ContentEntity(await this.contentService.findOne(contentId));
  }

  @ApiCrudOperation(Operation.DELETE, "content", ContentEntity, true)
  @Delete("/:contentId")
  remove(
    @Param("orgname") orgname: string,
    @Param("contentId") contentId: string
  ) {
    return this.contentService.remove(orgname, contentId);
  }

  @ApiCrudOperation(Operation.UPDATE, "content", ContentEntity, true)
  @Patch("/:contentId")
  async update(
    @Param("orgname") orgname: string,
    @Param("contentId") contentId: string,
    @Body() updateContentDto: UpdateContentDto
  ) {
    return new ContentEntity(
      await this.contentService.update(orgname, contentId, updateContentDto)
    );
  }
}
