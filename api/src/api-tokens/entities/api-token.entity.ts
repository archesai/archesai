import { ApiProperty } from "@nestjs/swagger";
import { ApiToken as _PrismaApiToken, RoleType } from "@prisma/client";
import { IsEnum, IsString } from "class-validator";

import { BaseEntity } from "../../common/dto/base.entity.dto";

export type ApiTokenModel = _PrismaApiToken;

export class ApiTokenEntity extends BaseEntity implements ApiTokenModel {
  @ApiProperty({
    default: "*",
    description: "The domains that can access this API token",
    example: "archesai.com,localhost:3000",
  })
  @IsString()
  domains: string;

  @ApiProperty({
    description: "The API token key. This will only be shown once",
    example: "********1234567890",
  })
  key: string;

  @ApiProperty({
    description: "The name of the API token",
    example: "My Token",
  })
  @IsString()
  name: string;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  orgname: string;

  @ApiProperty({
    description: "The role of the API token",
    enum: RoleType,
  })
  @IsEnum(RoleType)
  role: RoleType;

  @ApiProperty({
    description: "The username of the user who owns this API token",
    example: "jonathan",
  })
  username: string;

  constructor(apiToken: ApiTokenModel) {
    super();
    Object.assign(this, apiToken);
  }
}
