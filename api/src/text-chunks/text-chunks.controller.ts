import { Controller, Get, Param, Query } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { PaginatedDto } from "../common/paginated.dto";
import { TextChunkQueryDto } from "./dto/text-chunk-query.dto";
import { TextChunkEntity } from "./entities/text-chunk.entity";
import { TextChunksService } from "./text-chunks.service";

@ApiBearerAuth()
@ApiTags("Content - Vector Records")
@Controller("organizations/:orgname/content/:contentId/text-chunks")
export class TextChunksController
  implements
    BaseController<TextChunkEntity, undefined, TextChunkQueryDto, undefined>
{
  constructor(private readonly textChunksService: TextChunksService) {}

  @ApiCrudOperation(Operation.FIND_ALL, "vector record", TextChunkEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() textChunkQueryDto: TextChunkQueryDto,
    @Param("contentId") contentId?: string
  ) {
    const { count, results } = await this.textChunksService.findAll(
      orgname,
      textChunkQueryDto,
      contentId
    );
    const textChunkEntities = await Promise.all(
      results.map(async (textChunk) => {
        return new TextChunkEntity(textChunk);
      })
    );
    return new PaginatedDto<TextChunkEntity>({
      metadata: {
        limit: textChunkQueryDto.limit,
        offset: textChunkQueryDto.offset,
        totalResults: count,
      },
      results: textChunkEntities,
    });
  }

  @ApiCrudOperation(Operation.GET, "vector record", TextChunkEntity, true)
  @Get("/:textChunkId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("contentId") contentId: string
  ) {
    return new TextChunkEntity(await this.textChunksService.findOne(contentId));
  }
}
