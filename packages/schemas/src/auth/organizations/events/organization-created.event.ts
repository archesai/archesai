import type { OrganizationEntity } from '#auth/organizations/entities/organization.entity'
import type { UserEntity } from '#auth/users/entities/user.entity'

export class OrganizationCreatedEvent {
  public creator?: UserEntity
  public organization: OrganizationEntity

  constructor(event: {
    creator?: UserEntity
    organization: OrganizationEntity
  }) {
    this.organization = event.organization
    if (event.creator) {
      this.creator = event.creator
    }
  }
}
