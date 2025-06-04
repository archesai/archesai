import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { ORGANIZATION_ENTITY_KEY, OrganizationEntity } from '@archesai/domain'

/**
 * Repository for handling organizations.
 */
export class OrganizationRepository extends BaseRepository<OrganizationEntity> {
  constructor(databaseService: DatabaseService<OrganizationEntity>) {
    super(databaseService, ORGANIZATION_ENTITY_KEY, OrganizationEntity)
  }

  public async addOrRemoveCredits(
    orgname: string,
    numCredits: number
  ): Promise<OrganizationEntity> {
    const organization = await this.findByOrgname(orgname)
    const model = await this.update(ORGANIZATION_ENTITY_KEY, {
      billingEmail: organization.billingEmail,
      credits:
        numCredits < 0 ?
          organization.credits + numCredits
        : organization.credits - -1 * numCredits
    })
    return this.toEntity(model)
  }

  public async findByOrgname(orgname: string): Promise<OrganizationEntity> {
    const model = await this.findFirst({
      filter: {
        orgname: {
          equals: orgname
        }
      }
    })
    return this.toEntity(model)
  }
}
