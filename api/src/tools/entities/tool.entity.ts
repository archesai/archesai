import { ApiProperty } from "@nestjs/swagger";
import { Tool, ToolIOType } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEnum, IsOptional, IsString } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";

@Exclude()
export class ToolEntity extends BaseEntity implements Tool {
  @ApiProperty({
    description: "The tool description",
    example: "This tool converts a file to text, regardless of the file type.",
  })
  @Expose()
  @IsString()
  description: string;

  @ApiProperty({
    description: "The tools input type",
    enum: ToolIOType,
    example: "FILE",
  })
  @Expose()
  @IsEnum(ToolIOType)
  inputType: ToolIOType;

  @ApiProperty({
    description: "The tool's name",
    example: "extract-text",
  })
  @Expose()
  @IsString()
  name: string;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  @Expose()
  @IsString()
  orgname: string;

  @ApiProperty({
    description: "The tools output type",
    enum: ToolIOType,
    example: "TEXT",
  })
  @Expose()
  @IsEnum(ToolIOType)
  outputType: ToolIOType;

  @ApiProperty({
    description: "The tools output text",
    example: "Hello, world!",
  })
  @Expose()
  @IsOptional()
  text: string;

  @ApiProperty({ example: "https://example.com/example.mp4" })
  @Expose()
  @IsString()
  url: string;

  constructor(tool: Tool) {
    super();
    Object.assign(this, tool);
  }
}
