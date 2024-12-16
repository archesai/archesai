import Stripe from 'stripe'

export class Address {
  /**
   * City/District/Suburb/Town/Village.
   * @example 'San Francisco'
   */
  city: string | null

  /**
   * Two-letter country code (ISO 3166-1 alpha-2).
   * @example 'US'
   */
  country: string | null

  /**
   * Address line 1 (e.g., street, PO Box, or company name).
   * @example '123 Main Street'
   */
  line1: string | null

  /**
   * Address line 2 (e.g., apartment, suite, unit, or building).
   * @example 'Apt 4B'
   */
  line2: string | null

  /**
   * ZIP or postal code.
   * @example '94111'
   */
  postal_code: string | null

  /**
   * State/County/Province/Region.
   * @example 'CA'
   */
  state: string | null

  constructor(address: Stripe.Address | null) {
    this.city = address?.city || null
    this.country = address?.country || null
    this.line1 = address?.line1 || null
    this.line2 = address?.line2 || null
    this.postal_code = address?.postal_code || null
    this.state = address?.state || null
  }
}

export class BillingDetails {
  /**
   * Billing address associated with the payment method.
   */
  address: Address

  /**
   * Email address associated with the payment method.
   * @example 'customer@example.com'
   */
  email: string | null

  /**
   * Full name associated with the payment method.
   * @example 'Jane Doe'
   */
  name: string | null

  /**
   * Phone number associated with the payment method.
   * @example '+1234567890'
   */
  phone: string | null

  constructor(billingDetails: Stripe.PaymentMethod.BillingDetails) {
    this.address = new Address(billingDetails.address)
    this.email = billingDetails.email || null
    this.name = billingDetails.name || null
    this.phone = billingDetails.phone || null
  }
}

export class CardDetails {
  /**
   * Card brand (e.g., Visa, MasterCard).
   * @example 'visa'
   */
  brand: string

  /**
   * Two-letter ISO code representing the country of the card.
   * @example 'US'
   */
  country: string

  /**
   * Two-digit number representing the card’s expiration month.
   * @example 12
   */
  exp_month: number

  /**
   * Four-digit number representing the card’s expiration year.
   * @example 2024
   */
  exp_year: number

  /**
   * Unencrypted PAN tokens (optional, sensitive).
   */
  fingerprint: string | null

  /**
   * Card funding type (credit, debit, prepaid, unknown).
   * @example 'credit'
   */
  funding: string

  /**
   * The last four digits of the card.
   * @example '4242'
   */
  last4: string

  constructor(card: Stripe.PaymentMethod.Card) {
    this.brand = card.brand
    this.country = card.country
    this.exp_month = card.exp_month
    this.exp_year = card.exp_year
    this.funding = card.funding
    this.last4 = card.last4
    this.fingerprint = card.fingerprint || null
  }
}

export class PaymentMethodEntity {
  /**
   * Billing details associated with the payment method.
   */
  billing_details: BillingDetails

  /**
   * If the PaymentMethod is a card, this contains the card details.
   */
  card: CardDetails | null

  /**
   * ID of the customer this payment method is saved to.
   * @example 'cus_J0a1b2c3d4e5f6g7h8i9'
   */
  customer: string | null

  /**
   * Unique identifier for the payment method.
   * @example 'pm_1J2Y3A4B5C6D7E8F9G0H'
   */
  id: string

  /**
   * The type of the PaymentMethod. An example value is "card".
   * @example 'card'
   */
  type: string

  constructor(paymentMethod: Stripe.PaymentMethod) {
    this.id = paymentMethod.id
    this.type = paymentMethod.type
    this.customer =
      typeof paymentMethod.customer === 'string' ? paymentMethod.customer : null
    this.billing_details = new BillingDetails(paymentMethod.billing_details)
    this.card = paymentMethod.card ? new CardDetails(paymentMethod.card) : null
  }
}
