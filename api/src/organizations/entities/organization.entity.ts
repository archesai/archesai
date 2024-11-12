import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Organization as _PrismaOrganization, PlanType } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEmail, IsEnum, IsNumber } from "class-validator";

import { BaseEntity } from "../../common/entities/base.entity";

export type OrganizationModel = _PrismaOrganization;

@Exclude()
export class OrganizationEntity
  extends BaseEntity
  implements OrganizationModel
{
  @Expose()
  @ApiProperty({
    description: "The billing email to use for the organization",
    example: "example@test.com",
  })
  @IsEmail()
  billingEmail: string;

  @Expose()
  @ApiProperty({
    description:
      "The number of credits you have remaining for this organization",
    example: 500000,
  })
  @IsNumber()
  credits: number;

  @Expose()
  @ApiProperty({
    description: "The name of the organization to create",
    example: "organization-name",
  })
  orgname: string;

  @Expose()
  @ApiProperty({
    description: "The plan that the organization is subscribed to",
    enum: PlanType,
    example: "FREE",
  })
  @IsEnum(PlanType)
  plan: PlanType;

  @ApiHideProperty()
  stripeCustomerId: string;

  constructor(organization: OrganizationModel) {
    super();
    Object.assign(this, organization);
  }
}
