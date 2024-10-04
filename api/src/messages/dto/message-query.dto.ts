import { PickType } from "@nestjs/swagger";
import { ApiProperty } from "@nestjs/swagger";
import { IsEnum, IsOptional } from "class-validator";

import { SearchQueryDto } from "../../common/search-query";

export enum SortByField {
  CREATED = "createdAt",
}

export class MessageQueryDto extends PickType(SearchQueryDto, [
  "limit",
  "offset",
  "sortDirection",
] as const) {
  @ApiProperty({ default: "createdAt", enum: SortByField, required: false })
  @IsEnum(SortByField, { always: false })
  @IsOptional()
  sortBy? = "createdAt" as SortByField.CREATED;
}
