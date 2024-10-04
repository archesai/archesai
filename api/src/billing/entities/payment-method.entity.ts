import { ApiProperty } from "@nestjs/swagger";
import Stripe from "stripe";

export class Address {
  @ApiProperty({
    description: "City/District/Suburb/Town/Village.",
    example: "San Francisco",
    nullable: true,
  })
  city: null | string;

  @ApiProperty({
    description: "Two-letter country code (ISO 3166-1 alpha-2).",
    example: "US",
    nullable: true,
  })
  country: null | string;

  @ApiProperty({
    description: "Address line 1 (e.g., street, PO Box, or company name).",
    example: "123 Main Street",
    nullable: true,
  })
  line1: null | string;

  @ApiProperty({
    description: "Address line 2 (e.g., apartment, suite, unit, or building).",
    example: "Apt 4B",
    nullable: true,
  })
  line2: null | string;

  @ApiProperty({
    description: "ZIP or postal code.",
    example: "94111",
    nullable: true,
  })
  postal_code: null | string;

  @ApiProperty({
    description: "State/County/Province/Region.",
    example: "CA",
    nullable: true,
  })
  state: null | string;

  constructor(address: null | Stripe.Address) {
    this.city = address?.city || null;
    this.country = address?.country || null;
    this.line1 = address?.line1 || null;
    this.line2 = address?.line2 || null;
    this.postal_code = address?.postal_code || null;
    this.state = address?.state || null;
  }
}

export class BillingDetails {
  @ApiProperty({
    description: "Billing address associated with the payment method.",
    type: () => Address,
  })
  address: Address;

  @ApiProperty({
    description: "Email address associated with the payment method.",
    example: "customer@example.com",
    nullable: true,
  })
  email: null | string;

  @ApiProperty({
    description: "Full name associated with the payment method.",
    example: "Jane Doe",
    nullable: true,
  })
  name: null | string;

  @ApiProperty({
    description: "Phone number associated with the payment method.",
    example: "+1234567890",
    nullable: true,
  })
  phone: null | string;

  constructor(billingDetails: Stripe.PaymentMethod.BillingDetails) {
    this.address = new Address(billingDetails.address);
    this.email = billingDetails.email || null;
    this.name = billingDetails.name || null;
    this.phone = billingDetails.phone || null;
  }
}

export class CardDetails {
  @ApiProperty({
    description: "Card brand (e.g., Visa, MasterCard).",
    example: "visa",
  })
  brand: string;

  @ApiProperty({
    description: "Two-letter ISO code representing the country of the card.",
    example: "US",
  })
  country: string;

  @ApiProperty({
    description: "Two-digit number representing the card’s expiration month.",
    example: 12,
  })
  exp_month: number;

  @ApiProperty({
    description: "Four-digit number representing the card’s expiration year.",
    example: 2024,
  })
  exp_year: number;

  @ApiProperty({
    description: "Unencrypted PAN tokens (optional, sensitive).",
    nullable: true,
  })
  fingerprint: null | string;

  @ApiProperty({
    description: "Card funding type (credit, debit, prepaid, unknown).",
    example: "credit",
  })
  funding: string;

  @ApiProperty({
    description: "The last four digits of the card.",
    example: "4242",
  })
  last4: string;

  constructor(card: Stripe.PaymentMethod.Card) {
    this.brand = card.brand;
    this.country = card.country;
    this.exp_month = card.exp_month;
    this.exp_year = card.exp_year;
    this.funding = card.funding;
    this.last4 = card.last4;
    this.fingerprint = card.fingerprint || null;
  }
}

export class PaymentMethodEntity {
  @ApiProperty({
    description: "Billing details associated with the payment method.",
    type: () => BillingDetails,
  })
  billing_details: BillingDetails;

  @ApiProperty({
    description:
      "If the PaymentMethod is a card, this contains the card details.",
    nullable: true,
    type: () => CardDetails,
  })
  card: CardDetails | null;

  @ApiProperty({
    description: "ID of the customer this payment method is saved to.",
    example: "cus_J0a1b2c3d4e5f6g7h8i9",
    nullable: true,
  })
  customer: null | string;

  @ApiProperty({
    description: "Unique identifier for the payment method",
    example: "pm_1J2Y3A4B5C6D7E8F9G0H",
  })
  id: string;

  @ApiProperty({
    description: 'The type of the PaymentMethod. An example value is "card".',
    example: "card",
  })
  type: string;

  constructor(paymentMethod: Stripe.PaymentMethod) {
    this.id = paymentMethod.id;
    this.type = paymentMethod.type;
    this.customer =
      typeof paymentMethod.customer === "string"
        ? paymentMethod.customer
        : null;
    this.billing_details = new BillingDetails(paymentMethod.billing_details);
    this.card = paymentMethod.card ? new CardDetails(paymentMethod.card) : null;
  }
}
