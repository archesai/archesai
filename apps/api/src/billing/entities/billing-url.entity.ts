export class BillingUrlEntity {
  /**
   * The url that will bring you to the necessary stripe page
   * @example 'www.stripe.com/checkout/filchat-io'
   */
  url: string

  constructor(val: { url: string }) {
    this.url = val.url
  }
}
