import { PlanTypeEnum } from '@/src/organizations/entities/organization.entity'
import { Expose } from 'class-transformer'
import {
  IsEnum,
  IsNumber,
  IsObject,
  IsOptional,
  IsString,
  ValidateNested
} from 'class-validator'
import Stripe from 'stripe'

export class PlanMetadata {
  /**
   * The key of the metadata
   * @example 'STANDARD'
   */
  @IsOptional()
  @IsEnum(PlanTypeEnum)
  @Expose()
  key?: PlanTypeEnum
}

export class PlanEntity {
  /**
   * The currency of the plan
   * @example 'usd'
   */
  @IsString()
  @Expose()
  currency: string

  /**
   * The description of the plan
   * @example 'A plan for a small business'
   */
  @IsOptional()
  @IsString()
  @Expose()
  description: null | string

  /**
   * The ID of the plan
   * @example 'prod_1234567890'
   */
  @IsString()
  @Expose()
  id: string

  /**
   * The metadata of the plan
   */
  @ValidateNested()
  @Expose()
  metadata: PlanMetadata

  /**
   * The name of the plan
   * @example 'Small Business Plan'
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The ID of the price associated with the plan
   * @example 'price_1234567890'
   */
  @IsString()
  @Expose()
  priceId: string

  /**
   * The metadata of the price associated with the plan
   * @example { 'key': 'value' }
   */
  @IsObject()
  @Expose()
  priceMetadata: Record<string, string>

  /**
   * The interval of the plan
   */
  @IsOptional()
  @IsObject()
  @Expose()
  recurring: null | Stripe.Price.Recurring

  /**
   * The amount in cents to be charged on the interval specified
   * @example 1000
   */
  @IsNumber()
  @Expose()
  unitAmount: number

  constructor(partial: any) {
    Object.assign(this, partial)
  }
}
