import { ApiProperty } from "@nestjs/swagger";
import { ContentType } from "@prisma/client";
import { IsEnum, IsOptional } from "class-validator";

import { SearchQueryDto } from "../../common/search-query";

export enum SortByField {
  CREATED = "createdAt",
  USERNAME = "name",
}

export class ContentQueryDto extends SearchQueryDto {
  @ApiProperty({ default: "createdAt", enum: SortByField, required: false })
  @IsEnum(SortByField, { always: false })
  @IsOptional()
  sortBy? = "createdAt" as SortByField.CREATED;

  @ApiProperty({ enum: ContentType, required: false })
  @IsEnum(ContentType, { always: false })
  @IsOptional()
  type?: ContentType;
}
