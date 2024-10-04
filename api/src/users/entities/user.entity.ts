import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Member, Organization, User } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEmail, IsInt, IsString, MinLength } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";
import { MemberEntity } from "../../members/entities/member.entity";

export class UserEntity extends BaseEntity implements User {
  @Exclude()
  @ApiHideProperty()
  deactivated!: boolean;

  @Expose()
  @ApiProperty({
    description: "The user's default organization name",
    example: "my-organization",
  })
  defaultOrg: string;

  @ApiProperty({
    description: "The user's display name",
    example: "John Smith",
  })
  @Expose()
  displayName: string;

  @ApiProperty({
    description: "The user's e-mail",
    example: "example@archesai.com",
  })
  @IsEmail()
  @Expose()
  email!: string;

  @ApiProperty({
    description: "Whether or not the user's e-mail has been verified",
  })
  @Expose()
  emailVerified!: boolean;

  @ApiProperty({
    description: "The user's first name",
    example: "John",
  })
  @Expose()
  firstName: string;

  @ApiProperty({
    description: "The user's last name",
    example: "Smith",
  })
  @Expose()
  lastName: string;

  @Expose()
  @ApiProperty({
    description: "The memberships of the currently signed in user",
    type: [MemberEntity],
  })
  memberships: MemberEntity[];

  @Exclude()
  @ApiHideProperty()
  organizations?: Organization[];

  @Exclude()
  @ApiHideProperty()
  password: string;

  @ApiProperty({
    description: "The user's photo url",
    example: "/avatar.png",
  })
  @IsString()
  @Expose()
  photoUrl!: string;

  @Exclude()
  @ApiHideProperty()
  refreshToken: string;

  // Private Properties
  @Exclude()
  @ApiHideProperty()
  uid!: string;

  // Exposed Properties
  @ApiProperty({
    description: "The user's username",
    example: "jonathan",
    minLength: 5,
  })
  @MinLength(5)
  @Expose()
  username!: string;

  constructor(user: { memberships: Member[] } & User) {
    super();
    Object.assign(this, user);
    this.memberships = (this.memberships || []).map(
      (membership) => new MemberEntity(membership)
    );
    this.displayName = this.firstName
      ? this.firstName + " " + this.lastName
      : this.username;
  }
}
