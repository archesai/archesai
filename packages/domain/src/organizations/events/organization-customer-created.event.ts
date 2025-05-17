export class OrganizationCustomerCreatedEvent {
  public customer: string
  public orgname: string

  constructor(event: { customer: string; orgname: string }) {
    this.customer = event.customer
    this.orgname = event.orgname
  }
}
