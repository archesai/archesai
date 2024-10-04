import { ApiProperty } from "@nestjs/swagger";
import { IsBoolean, IsOptional, IsString } from "class-validator";

export class PathDto {
  @ApiProperty({
    default: false,
    description: "Whether or not this path points to a directory",
    example: false,
    required: false,
  })
  @IsOptional()
  @IsBoolean()
  isDir?: boolean = false;

  @ApiProperty({
    description: "The path that the file should upload to",
    example: "/location/in/storage",
  })
  @IsString()
  path: string;
}
