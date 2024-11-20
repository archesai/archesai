import { ApiProperty } from "@nestjs/swagger";
import { Tool as _PrismaTool, ToolIOType } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEnum, IsString } from "class-validator";

import { BaseEntity } from "../../common/entities/base.entity";

export type ToolModel = _PrismaTool;

@Exclude()
export class ToolEntity extends BaseEntity implements ToolModel {
  /**
   * The tool description
   * @example This tool converts a file to text, regardless of the file type.
   */
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
    description: "The tool's base path",
    example: "extract-text",
  })
  @Expose()
  @IsString()
  toolBase: string;

  constructor(tool: ToolModel) {
    super();
    Object.assign(this, tool);
  }
}
