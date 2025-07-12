export class OrganizationCustomerCreatedEvent {
  public customer: string
  public organizationId: string

  constructor(event: { customer: string; organizationId: string }) {
    this.customer = event.customer
    this.organizationId = event.organizationId
  }
}
