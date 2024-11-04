import { Controller, Get, Param, Query } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { TextChunkQueryDto } from "./dto/text-chunk-query.dto";
import { TextChunkEntity } from "./entities/text-chunk.entity";
import { TextChunksService } from "./text-chunks.service";

@ApiBearerAuth()
@ApiTags("Content - Text Chunks")
@Controller("organizations/:orgname/content/:contentId/text-chunks")
export class TextChunksController
  implements
    BaseController<TextChunkEntity, undefined, TextChunkQueryDto, undefined>
{
  constructor(private readonly textChunksService: TextChunksService) {}

  @ApiCrudOperation(Operation.FIND_ALL, "text chunk", TextChunkEntity, true)
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() textChunkQueryDto: TextChunkQueryDto,
    @Param("contentId") contentId?: string
  ) {
    return this.textChunksService.findAll(
      orgname,
      textChunkQueryDto,
      contentId
    );
  }

  @ApiCrudOperation(Operation.GET, "text chunk", TextChunkEntity, true)
  @Get("/:textChunkId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("contentId") contentId: string,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    @Param("textChunkId") textChunkId?: string
  ) {
    return this.textChunksService.findOne(contentId);
  }
}
