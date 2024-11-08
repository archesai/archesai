import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Member, Organization, PlanType } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEmail, IsEnum, IsNotEmpty, IsNumber } from "class-validator";

import { BaseEntity } from "../../common/dto/base.entity.dto";

@Exclude()
export class OrganizationEntity extends BaseEntity implements Organization {
  @Expose()
  @ApiProperty({
    description: "The billing email to use for the organization",
    example: "example@test.com",
  })
  @IsEmail()
  billingEmail!: string;

  @Expose()
  @ApiProperty({
    description:
      "The number of credits you have remaining for this organization",
    example: 500000,
  })
  @IsNumber()
  credits!: number;

  @ApiHideProperty()
  members: Member[];

  // Exposed Properties
  @Expose()
  @ApiProperty({
    description: "The name of the organization to create",
    example: "organization-name",
  })
  @IsNotEmpty()
  orgname!: string;

  @Expose()
  @ApiProperty({
    description: "The plan that the organization is subscribed to",
    enum: PlanType,
    example: "FREE",
  })
  @IsEnum(PlanType)
  plan!: PlanType;

  @ApiHideProperty()
  stripeCustomerId!: string;

  constructor(organization: Partial<Organization>) {
    super();
    Object.assign(this, organization);
  }
}
