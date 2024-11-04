import { ApiProperty, PickType } from "@nestjs/swagger";
import { ApiToken, RoleType } from "@prisma/client";
import { IsEnum, IsString, ValidateNested } from "class-validator";

import { ChatbotEntity } from "../../chatbots/entities/chatbot.entity";
import { BaseEntity } from "../../common/base-entity.dto";

export class ChatbotsFieldItem extends PickType(ChatbotEntity, [
  "id",
  "name",
]) {}

export class ApiTokenEntity extends BaseEntity implements ApiToken {
  @ApiProperty({
    description: "The chatbots this API token has access to",
    example: [
      { id: "uuid-uuid-uuid-uuid", name: "Arches API Documentation Chatbot" },
    ],
    type: [ChatbotsFieldItem],
  })
  @ValidateNested({ each: true })
  chatbots: ChatbotsFieldItem[];

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
