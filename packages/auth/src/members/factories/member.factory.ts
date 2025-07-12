import { faker } from '@faker-js/faker'

import type { MemberEntity } from '@archesai/schemas'

import { RoleTypes } from '@archesai/schemas'

export function createRandomMember(
  overrides?: Partial<MemberEntity>
): MemberEntity {
  return {
    createdAt: faker.date.recent().toISOString(),
    id: faker.string.uuid(),
    invitationId: faker.string.uuid(),
    organizationId: faker.internet.domainName(),
    role: faker.helpers.arrayElement(RoleTypes),
    updatedAt: faker.date.recent().toISOString(),
    userId: faker.string.uuid(),
    ...overrides
  }
}
