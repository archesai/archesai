import { ApiHideProperty, ApiProperty, PickType } from "@nestjs/swagger";
import { ApiToken, Organization, RoleType, User } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEnum, IsOptional, IsString, ValidateNested } from "class-validator";

import { ChatbotEntity } from "../../chatbots/entities/chatbot.entity";
import { BaseEntity } from "../../common/base-entity.dto";

export class AgentsFieldItem extends PickType(ChatbotEntity, ["id", "name"]) {}

export class ApiTokenEntity extends BaseEntity implements ApiToken {
  @ApiProperty({
    description: "The chatbots this API token has access to",
    example: [
      { id: "uuid-uuid-uuid-uuid", name: "Arches API Documentation Agent" },
    ],
    type: [AgentsFieldItem],
  })
  @Expose()
  @IsOptional()
  @ValidateNested({ each: true })
  chatbots: AgentsFieldItem[];

  @ApiProperty({
    default: "*",
    description: "The domains that can access this API token",
    example: "archesai.com,localhost:3000",
    required: false,
  })
  @Expose()
  @IsString()
  domains: string = "*";

  @ApiProperty({
    description: "The API token key. This will only be shown once",
    example: "********1234567890",
  })
  @Expose()
  key: string;

  @ApiProperty({
    description: "The name of the API token",
    example: "My Token",
  })
  @Expose()
  @IsString()
  name: string;

  // Private Properties
  @ApiHideProperty()
  @Exclude()
  organization: Organization;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  @Expose()
  orgname: string;

  // Public Properties
  @ApiProperty({ description: "The role of the API token", enum: RoleType })
  @Expose()
  @IsEnum(RoleType)
  role: RoleType;

  @ApiHideProperty()
  @Exclude()
  user: User;

  @ApiProperty({
    description: "The username of the user who owns this API token",
    example: "jonathan",
  })
  @Expose()
  username: string;

  constructor(
    token: {
      chatbots: { id: string; name: string }[];
    } & ApiToken
  ) {
    super();
    Object.assign(this, token);
    this.chatbots = token.chatbots.map((a) => ({ id: a.id, name: a.name }));
  }
}
