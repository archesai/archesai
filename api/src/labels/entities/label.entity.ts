import { ApiProperty } from "@nestjs/swagger";
import { Label as _PrismaLabel } from "@prisma/client";
import { Expose } from "class-transformer";
import { IsOptional } from "class-validator";

import { BaseEntity } from "../../common/entities/base.entity";

export type LabelModel = _PrismaLabel;

export class LabelEntity extends BaseEntity implements LabelModel {
  @ApiProperty({
    default: "New Chat",
    description: "The chat label name",
    example: "What are the morals of the story in Aesop's Fables?",
    required: false,
  })
  @Expose()
  @IsOptional()
  name: string;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  @Expose()
  orgname: string;

  constructor(label: LabelModel) {
    super();
    Object.assign(this, label);
  }
}
