import { SearchQueryDto } from "@/src/common/search-query";
import { IntersectionType, PickType } from "@nestjs/swagger";

import { RunEntity } from "../entities/run.entity";

export class RunQueryDto extends IntersectionType(
  SearchQueryDto,
  PickType(RunEntity, ["toolId"] as const)
) {}
