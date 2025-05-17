export class OrganizationCustomerSubscriptionUpdatedEvent {
  public credits?: number
  public customer: string
  public orgname: string
  public planType?: string

  constructor(event: {
    credits?: number
    customer: string
    orgname: string
    planType: string
  }) {
    this.customer = event.customer
    this.orgname = event.orgname
    if (event.credits) {
      this.credits = event.credits
    }
    this.planType = event.planType
  }
}
