export class OrganizationCustomerSubscriptionUpdatedEvent {
  public credits?: number
  public customer: string
  public organizationId: string
  public planType?: string

  constructor(event: {
    credits?: number
    customer: string
    organizationId: string
    planType: string
  }) {
    this.customer = event.customer
    this.organizationId = event.organizationId
    if (event.credits) {
      this.credits = event.credits
    }
    this.planType = event.planType
  }
}
