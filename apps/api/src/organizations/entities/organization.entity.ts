import { ApiHideProperty } from '@nestjs/swagger'
import { Organization as _PrismaOrganization } from '@prisma/client'
import { Expose } from 'class-transformer'
import { IsEmail, IsEnum, IsNumber, IsString } from 'class-validator'

import { BaseEntity } from '../../common/entities/base.entity'

export type OrganizationModel = _PrismaOrganization

export enum PlanTypeEnum {
  FREE = 'FREE',
  BASIC = 'BASIC',
  STANDARD = 'STANDARD',
  PREMIUM = 'PREMIUM',
  UNLIMITED = 'UNLIMITED'
}

export class OrganizationEntity
  extends BaseEntity
  implements OrganizationModel
{
  /**
   * The billing email to use for the organization
   * @example 'example@test.com'
   */
  @IsEmail()
  @Expose()
  billingEmail: string

  /**
   * The number of credits you have remaining for this organization
   * @example 500000
   */
  @IsNumber()
  @Expose()
  credits: number

  /**
   * The name of the organization
   * @example 'organization-name'
   */
  @IsString()
  @Expose()
  orgname: string

  /**
   * The plan that the organization is subscribed to
   * @example FREE
   */
  @IsEnum(PlanTypeEnum)
  @Expose()
  plan: PlanTypeEnum

  @ApiHideProperty()
  stripeCustomerId: string

  constructor(organization: OrganizationModel) {
    super()
    Object.assign(this, organization)
  }
}
