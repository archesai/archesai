import { ApiProperty } from "@nestjs/swagger";
import { IsEnum, IsOptional } from "class-validator";

import { SearchQueryDto } from "../../common/search-query";

export enum SortByField {
  CREATED = "createdAt",
  USERNAME = "username",
}

export class MemberQueryDto extends SearchQueryDto {
  @ApiProperty({ default: "username", enum: SortByField, required: false })
  @IsEnum(SortByField, { always: false })
  @IsOptional()
  sortBy? = "username" as SortByField.USERNAME;
}
