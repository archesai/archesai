import { Controller, Get, Param, Query } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { PaginatedDto } from "../common/paginated.dto";
import { VectorRecordQueryDto } from "./dto/vector-record-query.dto";
import { VectorRecordEntity } from "./entities/vector-record.entity";
import { VectorRecordService } from "./vector-record.service";

@ApiBearerAuth()
@ApiTags("Content - Vector Records")
@Controller("organizations/:orgname/content/:contentId/vector-records")
export class VectorRecordController
  implements
    BaseController<
      VectorRecordEntity,
      undefined,
      VectorRecordQueryDto,
      undefined
    >
{
  constructor(private readonly vectorRecordService: VectorRecordService) {}

  @ApiCrudOperation(
    Operation.FIND_ALL,
    "vector record",
    VectorRecordEntity,
    true
  )
  @Get("/")
  async findAll(
    @Param("orgname") orgname: string,
    @Query() vectorRecordQueryDto: VectorRecordQueryDto,
    @Param("contentId") contentId?: string
  ) {
    const { count, results } = await this.vectorRecordService.findAll(
      orgname,
      vectorRecordQueryDto,
      contentId
    );
    const vectorRecordEntities = await Promise.all(
      results.map(async (vectorRecord) => {
        return new VectorRecordEntity(vectorRecord);
      })
    );
    return new PaginatedDto<VectorRecordEntity>({
      metadata: {
        limit: vectorRecordQueryDto.limit,
        offset: vectorRecordQueryDto.offset,
        totalResults: count,
      },
      results: vectorRecordEntities,
    });
  }

  @ApiCrudOperation(Operation.GET, "vector record", VectorRecordEntity, true)
  @Get("/:vectorRecordId")
  async findOne(
    @Param("orgname") orgname: string,
    @Param("contentId") contentId: string
  ) {
    return new VectorRecordEntity(
      await this.vectorRecordService.findOne(contentId)
    );
  }
}
