import { ApiProperty } from "@nestjs/swagger";

export class BillingUrlEntity {
  @ApiProperty({
    description: "The url that will bring you to the necessary stripe page",
    example: "www.stripe.com/checkout/filchat-io",
  })
  url!: string;

  constructor(val) {
    this.url = val.url;
  }
}
