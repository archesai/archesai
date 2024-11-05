import { ApiProperty, IntersectionType, PickType } from "@nestjs/swagger";
import { IsBoolean, IsEnum, IsOptional } from "class-validator";

import { Granularity } from "../../common/aggregated-field.dto";
import { SearchQueryDto } from "../../common/search-query";
import { ThreadEntity } from "../entities/thread.entity";

export enum SortByField {
  CREATED = "createdAt",
  CREDITS = "credits",
}

export class ThreadQueryDto extends IntersectionType(
  SearchQueryDto,
  PickType(ThreadEntity, ["chatbotId"] as const)
) {
  @ApiProperty({
    default: undefined,
    description: "The granularity to use for ranged aggregates",
    enum: Granularity,
    required: false,
  })
  @IsOptional()
  @IsEnum(Granularity, { always: false })
  aggregateGranularity?: Granularity;

  @ApiProperty({
    default: false,
    description: "Whether or not to include aggregates in the response",
    required: false,
  })
  @IsBoolean({ always: false })
  @IsOptional()
  aggregates?: boolean = false;

  @ApiProperty({ default: "createdAt", enum: SortByField, required: false })
  @IsEnum(SortByField, { always: false })
  @IsOptional()
  sortBy? = "createdAt" as SortByField.CREATED;
}
