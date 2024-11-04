import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Member, Organization, RoleType, User } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEnum, IsNotEmpty } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";

export class MemberEntity extends BaseEntity implements Member {
  @ApiProperty({
    description: "Whether the invite was accepted",
    example: false,
  })
  @Expose()
  @IsNotEmpty()
  inviteAccepted: boolean;

  // Exposed Properties
  @ApiProperty({
    description: "The invited email of this member",
    example: "invited-user@archesai.com",
  })
  @Expose()
  @IsNotEmpty()
  inviteEmail: string;

  @ApiProperty({
    description: "The name of this member",
    example: "jonathan",
  })
  name: string;

  @ApiHideProperty()
  @Exclude()
  organization: Organization;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  @Expose()
  orgname: string;

  @ApiProperty({ description: "The role of the member", enum: RoleType })
  @Expose()
  @IsEnum(RoleType)
  role: RoleType;

  // Private Properties
  @ApiHideProperty()
  @Exclude()
  user: User;

  @ApiProperty({
    description: "The username of this member",
    example: "jonathan",
    required: false,
    type: String,
  })
  @Expose()
  username: null | string;

  constructor(member: Member) {
    super();
    Object.assign(this, member);
    this.name = this.username || "";
  }
}
