import { ApiProperty } from "@nestjs/swagger";
import Stripe from "stripe";

export class PlanEntity {
  @ApiProperty({
    description: "The currency of the plan",
    example: "usd",
  })
  currency: string;

  @ApiProperty({ example: "A plan for a small business", required: false })
  description: null | string;

  @ApiProperty({
    description: "The ID of the plan",
    example: "prod_1234567890",
  })
  id: string;

  @ApiProperty({ required: false })
  metadata: Record<string, string>;

  @ApiProperty({
    description: "The name of the plan",
    example: "Small Business Plan",
  })
  name: string;

  @ApiProperty({
    description: "The ID of the price associated with the plan",
    example: "price_1234567890",
  })
  priceId: string;

  @ApiProperty({ required: false })
  priceMetadata: Record<string, string>;

  @ApiProperty({ required: false })
  recurring: null | Stripe.Price.Recurring;

  @ApiProperty({
    description: "The amount in cents to be charged on the interval specified",
    example: 1000,
  })
  unitAmount: number;

  constructor(partial: Partial<PlanEntity>) {
    Object.assign(this, partial);
  }
}
