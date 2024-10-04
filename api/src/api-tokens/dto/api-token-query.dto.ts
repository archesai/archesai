import { ApiProperty } from "@nestjs/swagger";
import { IsOptional, IsString } from "class-validator";

import { SearchQueryDto } from "../../common/search-query";

export class ApiTokenQueryDto extends SearchQueryDto {
  @IsOptional()
  @IsString()
  @ApiProperty({
    description: "The name to search for",
    example: "API Token 1",
  })
  name?: string;
}
